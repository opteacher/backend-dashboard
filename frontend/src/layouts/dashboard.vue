<template>
<el-container class="h-100">
    <el-header class="p-0"><header-bar/></el-header>
    <el-aside class="sidebar" width="20vw"><side-bar/></el-aside>
    <el-main class="main p-0"><slot/></el-main>
</el-container>
</template>

<script>
import headerBar from "../components/headerBar"
import sideBar from "../components/sideBar"

export default {
    async created() {
    },
    components: {
        "header-bar": headerBar,
        "side-bar": sideBar
    }
}
</script>

<style lang="scss">
.sidebar {
    position: fixed;
    top: 60px;
    bottom: 0;
}
.main {
    position: fixed;
    top: 60px;
    bottom: 0;
    left: 20vw;
    right: 0;
    // background-color: #f8f9fa !important;
    background-color: cadetblue !important;
}
</style>