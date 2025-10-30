<script setup lang="ts">
import { RouterView } from "vue-router"
import Header from "./components/Header.vue"
import Loading from "./components/Loading.vue"
</script>

<template>
	<Header/>
	<RouterView v-slot="{Component}">
		<template v-if="Component">
			<Transition mode="out-in">
				<KeepAlive>
					<Suspense>
						<component :is="Component"></component>
						<template #fallback>
							<Loading/>
						</template>
					</Suspense>
				</KeepAlive>
			</Transition>
		</template>
	</RouterView>
</template>
