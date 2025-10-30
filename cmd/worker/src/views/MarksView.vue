<script setup lang="ts">
import type { Lesson, Mark } from "../client/types.ts"
import { getMarks } from "../client/client.ts"
import SlotMark from "../components/Slot.vue"
import { humanizeLesson } from "../utils.ts"

const data = await getMarks()

function showAverageForClass() {
	window.Telegram.WebApp.showAlert("Средний балл по классу")
}

function showAverageForStudent() {
	window.Telegram.WebApp.showAlert("Средний балл ученика")
}

function isPass(mark: Mark): boolean {
	return mark.value === "Б" || mark.value === "Н" || mark.value === "У"
}

function isEmptyLesson(lesson: Lesson): boolean {
	return !lesson.slots.filter(slot => data.showPasses || !isPass(slot.mark)).length
}
</script>

<template>
	<section class="marks">
		<div class="marks-header">Оценки за {{ data.studyPeriod.text }}</div>
		<div class="marks-lessons">
			<template v-for="lesson in data.lessons">
				<div v-if="data.showEmptyLessons || !isEmptyLesson(lesson)" class="lesson">
					<div class="lesson-header">
						<div class="lesson-name">{{ humanizeLesson(lesson.subject) }}</div>
						<div class="lesson-average">
							<div
								class="mark average-mark"
								v-if="lesson.averages.for_class"
								:data-mark=lesson.averages.for_class[1]
								@click="showAverageForClass">
								{{ lesson.averages.for_class[0] }}
							</div>
							<div class="mark average-mark"
							  v-if="lesson.averages.for_student"
								:data-mark=lesson.averages.for_student[1]
								@click="showAverageForStudent">
								{{ lesson.averages.for_student[0] }}
							</div>
						</div>
					</div>
					<div class="lesson-marks">
						<template v-for="slot in lesson.slots">
							<SlotMark v-if="data.showPasses || !isPass(slot.mark)" :slot/>
						</template>
					</div>
				</div>
			</template>
		</div>
	</section>
</template>

<style>
.mark {
	background-color: var(--tg-theme-section-bg-color);
	border-radius: 10px;
	display: flex;
	justify-content: center;
	align-items: center;
	cursor: pointer;
	-webkit-tap-highlight-color: transparent;
}

.mark[data-mark*="1"],
.mark[data-mark*="2"] {
	color: #ef4444;
	background-color: #dc262640
}

.mark[data-mark*="3"] {
	color: #f97316;
	background-color: #ea580c40;
}

.mark[data-mark*="4"] {
	color: #3b82f6;
	background-color: #2563eb40;
}

.mark[data-mark*="5"] {
	color: #22c55e;
	background-color: #16a34a40;
}
</style>

<style scoped>
.marks {
	padding: 10px;
}

.marks-header {
	text-align: center;
	font-weight: bold;
	padding-bottom: 10px;
}

.marks-lessons {
	width: 100%;
	max-width: 700px;
	border-radius: 10px;
	margin-bottom: 0;
	overflow: hidden;
}

.lesson {
	display: flex;
	flex-direction: column;
	padding: 9px 14px;
	background-color: var(--tg-theme-secondary-bg-color);
	margin-bottom: 1px;
}

.lesson-header {
	padding-bottom: 10px;
	padding-left: 3px;
	display: flex;
	justify-content: space-between;
	align-items: center;
	flex-wrap: wrap;
}

.lesson-average {
	display: flex;
	flex-wrap: wrap;
}

.lesson-name {
	font-weight: bold;
}

.lesson-marks {
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(40px, 1fr));
	gap: 2px;
}

.average-mark {
	padding: 1px 10px;
	margin-left: 2px
}
</style>
