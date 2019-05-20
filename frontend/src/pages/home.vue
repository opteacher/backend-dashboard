<template>
<dashboard>
    <tool-bar/>
    <model-card/>
</dashboard>
</template>

<script>
import dashboard from "../layouts/dashboard"
import toolBar from "../components/toolBar"
import modelCard from "../components/modelCard"

export default {
    async created() {
    },
    components: {
        "dashboard": dashboard,
        "model-card": modelCard,
        "tool-bar": toolBar,
    }
}
</script>


