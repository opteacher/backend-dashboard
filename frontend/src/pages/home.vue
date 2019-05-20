<template>
<dashboard>
    <tool-bar @addModel="addModel"/>
    <div>
        <model-card v-for="model in models" :key="model.id" :model="model"/>
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
            models: []
        }
    },
    methods: {
        addModel(model) {
            this.models.push(model)
        }
    },
    async created() {
    },
    components: {
        "dashboard": dashboard,
        "model-card": modelCard,
        "tool-bar": toolBar,
    }
}
</script>


