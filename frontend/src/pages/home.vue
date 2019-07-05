<template>
<dashboard>
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

