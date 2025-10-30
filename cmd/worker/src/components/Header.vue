<script setup lang="ts">
import { onMounted, ref } from "vue"
import { getStudent } from "../client/client"

const studentNameAcronym = ref("")

onMounted(async () => {
	studentNameAcronym.value = (await getStudent()).studentNameAcronym
})
</script>

<template>
	<section class="header" v-if="studentNameAcronym">
		{{ studentNameAcronym }}
	</section>
</template>

<style scoped>
.header {
	font-variant: all-small-caps;
	background-color: var(--tg-theme-button-color);
	color: var(--tg-theme-button-text-color);
	padding: 2px 0;
	height: 25px;
}
</style>
