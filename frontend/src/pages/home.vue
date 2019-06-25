<template>
<dashboard>
    <!-- <div class="models-panel">
        <model-card v-for="model in models" :key="model.id" :model="model" @delete-model="deleteModel"/>
        <relation-link v-for="link in relations" :key="link.id" :relation="link"
                :model1="findModel(link.model1)"
                :model2="findModel(link.model2)"/>
    </div> -->
    <models-panel/>
</dashboard>
</template>

<script>
import utils from "../utils"
import dashboard from "../layouts/dashboard"
import modelCard from "../components/modelCard"
import relationLink from "../components/relationLink"
import modelsPanel from "../panels/modelsPanel"

export default {
    components: {
        "dashboard": dashboard,
        "model-card": modelCard,
        "relation-link": relationLink,
        "models-panel": modelsPanel
    }
}
</script>

<style lang="scss">
.models-panel {
    position: relative;
    height: 100%;
}
</style>

