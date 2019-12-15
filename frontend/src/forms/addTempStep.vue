<template>
<el-form ref="add-temp-step-form" :model="tempStep" label-width="80px">
    <el-form-item label="标识">
        <el-input v-model="tempStep.key"/>
    </el-form-item>
    <el-form-item label="组">
        <el-input v-model="tempStep.group"/>
    </el-form-item>
    <el-form ref="add-require-form" :model="addRequireForm" label-width="80px">
        <el-form-item label="添加依赖" :rules="[
            {required: true, message: '请选择要添加的依赖', trigger: 'blur'}
        ]" prop="reqName">
            <el-col :span="20">
                <el-select class="w-100" v-model="addRequireForm.reqName">
                    <el-option v-for="req in requires" :key="req.label" :value="req.value">
                        {{req.label}}
                    </el-option>
                </el-select>
            </el-col>
            <el-col :span="4">
                <el-button class="float-right" @click="addRequire">添加</el-button>
            </el-col>
        </el-form-item>
    </el-form>
    <el-form-item v-if="tempStep.requires.length !== 0" label="依赖">
        <div class="interval-container">
            <el-tag class="interval-item" v-for="req in tempStep.requires" :key="req" closable @close="rmvRequire(req)">
                {{req}}
            </el-tag>
        </div>
    </el-form-item>
    <el-form-item label="描述">
        <el-input type="textarea" :rows="2" v-model="tempStep.desc"/>
    </el-form-item>
    <str-set-input label="输入" :dataSet="tempStep.inputs"/>
    <str-set-input label="输出" :dataSet="tempStep.outputs"/>
    <el-form-item label="代码">
        <el-input type="textarea" :rows="4" v-model="tempStep.code"/>
    </el-form-item>
</el-form>
</template>

<script>
import strSetInput from "../components/strSetInput"

export default {
    components: {
        "str-set-input": strSetInput
    },
    data() {
        return {
            tempStep: {
                key: "",
                group: "",
                requires: [],
                desc: "",
                inputs: [],
                outputs: [],
                code: ""
            },
            addRequireForm: {
                reqName: ""
            },
            requires: [{
                label: "json",
                value: "encoding/json"
            }]
        }
    },
    methods: {
        addRequire() {
            this.$refs["add-require-form"].validate(valid => {
                if (!valid) {
                    return
                }
                this.tempStep.requires.push(this.addRequireForm.reqName)
                this.addRequireForm.reqName = ""
            })
        },
        rmvRequire(rmvReq) {
            this.tempStep.requires.pop(req => req === rmvReq)
        }
    }
}
</script>