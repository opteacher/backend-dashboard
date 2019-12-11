<template>
<el-input type="textarea" :autosize="{minRows: 5, maxRows: 100}" readonly="true" v-model="code"/>
</template>

<script>
export default {
    props: {
        "code": String
    }
}
</script>