<template>
<el-card v-if="show" class="box-card" style="width: 30vw">
    <div slot="header" class="clearfix">
        <span class="card-name">卡片名称</span>
        <el-link class="float-right" @click="deleteModel">
            <i class="el-icon-close"/>
        </el-link>
    </div>
    <div class="text item">
        tttttt
    </div>
</el-card>
</template>

<script>
export default {
    data() {
        return {
            show: true
        }
    },
    async created() {

    },
    methods: {
        deleteModel() {
            this.show = false
        }
    }
}
</script>


<style lang="scss">
.card-name {
    cursor: move;
}
</style>
