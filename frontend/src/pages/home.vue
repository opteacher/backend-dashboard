<template>
<dashboard>
    <tool-bar @add-model="addModel"/>
    <div class="models-panel">
        <model-card v-for="model in models" :key="model.id" :model="model" @deleteModel="deleteModel"/>
    </div>
</dashboard>
</template>

<script>
import dashboard from "../layouts/dashboard"
import toolBar from "../components/toolBar"
import modelCard from "../components/modelCard"

export default {
    data() {
        return {
            movingModel: null,
            models: []
        }
    },
    methods: {
        addModel(model) {
            model.x = 0
            model.y = 0
            this.models.push(model)
        },
        deleteModel(modelID) {
            this.models.pop(ele => ele.id === modelID)
        }
    },
    components: {
        "dashboard": dashboard,
        "model-card": modelCard,
        "tool-bar": toolBar,
    }
}
</script>

<style lang="scss">
.models-panel {
    position: relative;
    height: 100%;
}
</style>

