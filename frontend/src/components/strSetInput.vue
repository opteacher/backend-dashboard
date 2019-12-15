<template>
<div>
    <el-form ref="add-str-form" :model="addForm" label-width="80px">
        <el-form-item :label="`添加${label}`" :rules="[
            {required: true, message: `不能添加空的${label}`, trigger: 'blur'}
        ]" prop="str">
            <el-col :span="20">
                <el-input v-model="addForm.str"/>
            </el-col>
            <el-col :span="4">
                <el-button class="float-right" @click="addStr">添加</el-button>
            </el-col>
        </el-form-item>
    </el-form>
    <el-form-item v-if="dataSet.length !== 0" :label="label">
        <div class="interval-container">
            <el-tag class="interval-item" v-for="str in dataSet" :key="str" closable @close="rmvStr(str)">
                {{str}}
            </el-tag>
        </div>
    </el-form-item>
</div>
</template>

<script>
export default {
    props: {
        "label": String,
        "dataSet": Array
    },
    data() {
        return {
            addForm: {
                str: ""
            }
        }
    },
    methods: {
        addStr() {
            this.$refs["add-str-form"].validate(valid => {
                if (!valid) {
                    return
                }
                this.dataSet.push(this.addForm.str)
                this.addForm.str = ""
            })
        },
        rmvStr(rmvStr) {
            this.dataSet.pop(str => str === rmvStr)
        }
    }
}
</script>