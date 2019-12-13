<template>
<div ref="chk-load-res">
    <p class="mb-10">检测到DAO实例有模板步骤和接口，是否导入？</p>
    <ul class="list-inline">
        <li class="list-inline-item">
            <el-checkbox v-model="loadTmpStep" :disabled="!hrefAvai(daoImpl.tmpStepHref)">模板步骤</el-checkbox>
        </li>
        <li class="list-inline-item">
            <el-checkbox v-model="loadApiTemp" :disabled="!hrefAvai(daoImpl.apiTempHref)">模板API</el-checkbox>
        </li>
    </ul>
</div>
</template>

<script>
export default {
    props: {
        "daoImpl": Object
    },
    data() {
        return {
            loadTmpStep: true,
            loadApiTemp: true
        }
    },
    methods: {
        hrefAvai(href) {
            return href && href.length !== 0
        }
    }
}
</script>