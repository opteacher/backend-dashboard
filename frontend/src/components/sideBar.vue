<template>
<el-menu :default-active="activeIdx" class="el-menu-vertical-demo h-100 br-0" @select="handleSelect" background-color="#545c64" text-color="#fff" active-text-color="#ffd04b">
    <el-menu-item index="models">
        <i class="el-icon-menu"></i>
        <span slot="title">Models</span>
    </el-menu-item>
    <el-menu-item index="services">
        <i class="el-icon-s-help"></i>
        <span slot="title">Services</span>
    </el-menu-item>
    <el-menu-item index="steps">
        <i class="el-icon-s-shop"></i>
        <span slot="title">OperSteps</span>
    </el-menu-item>
    <el-menu-item index="dao">
        <i class="el-icon-share"></i>
        <span slot="title">DAO</span>
    </el-menu-item>
    <el-menu-item index="exports">
        <i class="el-icon-s-promotion"></i>
        <span slot="title">Exports</span>
    </el-menu-item>
</el-menu>
</template>

<script>
export default {
    data() {
        return {
            activeIdx: "models"
        }
    },
    created() {
        this.activeIdx = this.$route.path.split("/").pop()
    },
    methods: {
        handleSelect(key) {
            this.$router.push(`/backend-dashboard/frontend/${key}`)
        }
    }
}
</script>
