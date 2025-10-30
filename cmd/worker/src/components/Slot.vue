<script setup lang="ts">
import type { PopupParams } from "telegram-web-app"
import type { Slot } from "../client/types.ts"

const { slot } = defineProps<{ slot: Slot }>()

const date = slot.lesson_date.slice(8, 10) + "." + slot.lesson_date.slice(5, 7)

function onClick(): void {
	const params: PopupParams = {
		title: new Date(slot.updated_at).toLocaleString("ru"),
		message: slot.text || "Нет комментария",
	}

	window.Telegram.WebApp.showPopup(params)
}
</script>

<template>
	<div class="slot">
		<div class="mark slot-mark" :data-mark=slot.mark.value @click="onClick">
			{{ slot.mark.value }}
		</div>
		<div class="slot-date">{{ date }}</div>
	</div>
</template>

<style scoped>
.slot {
	display: flex;
	flex-direction: column;
	flex-wrap: wrap;
	align-items: center;
}

.slot-mark {
	width: 40px;
	height: 40px;
	font-size: large;
}

.slot-date {
	font-size: 0.8em
}
</style>
