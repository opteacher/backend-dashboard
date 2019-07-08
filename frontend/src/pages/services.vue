<template>
<dashboard>
    <info-bar @sel-interface="selInterface"/>
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
    methods: {
        selInterface(selItf) {
            console.log(selItf)
        }
    }
}
</script>
