<template>
<el-collapse v-model="activeStruct" accordion>
    <el-collapse-item v-for="struct in structs" :key="struct.name" :title="struct.name" :name="struct.name">
        <div v-for="prop in struct.props" :key="prop.name">
            {{prop.name}}<el-tag>prop.type</el-tag>
        </div>
    </el-collapse-item>
</el-collapse>
</template>

<script>
export default {
    props: {
        "showFlag": Boolean
    },
    data() {
        return {
            structs: [],
            activeStruct: ""
        }
    },
    created() {

    },
    watch: {
        showFlag() {

        }
    },
    methods: {
        async refresh() {
            
        }
    }
}
</script>
