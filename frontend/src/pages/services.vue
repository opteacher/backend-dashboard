<template>
<dashboard>
    <info-bar/>
</dashboard>
</template>

<script>
import dashboard from "../layouts/dashboard"
import infoBar from "../components/infoBar"

export default {
    components: {
        "dashboard": dashboard,
        "info-bar": infoBar
    },
}
</script>
