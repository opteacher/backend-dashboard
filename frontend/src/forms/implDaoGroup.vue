<template>
<el-form ref="impl-dao-group-form" :model="selImpl">
    <el-form v-model="searchForm">
        <el-form-item :rules="[
            { required: true, message: '需要指定DAO实例名', trigger: 'blur' },
            { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
        ]" prop="search">
            <el-input placeholder="请输入DAO实例名" v-model="searchForm.search">
                <el-button slot="append" icon="el-icon-search"></el-button>
            </el-input>
        </el-form-item>
    </el-form>
    <ul class="list-unstyled infinite-list" v-infinite-scroll="loadImpls" style="height:30vh;overflow:auto">
        <li v-for="i in rowNum" :key="i" class="infinite-list-item">
            <el-row class="m-0 mb-20" :gutter="10">
                <el-col v-for="j in [0, 1, 2, 3]" :key="j" :span="6">
                    <el-link :underline="false" v-if="avaImpls[i + j]" @click="selImpl = avaImpls[i + j]">
                        <el-card :class="{active: selImpl === avaImpls[i + j]}" shadow="hover">
                            <div class="clearfix" slot="header">
                                <span>{{avaImpls[i + j].name}}</span>
                            </div>
                            <div style="font-size:0.2em">{{avaImpls[i + j].desc}}</div>
                        </el-card>
                    </el-link>
                </el-col>
            </el-row>
        </li>
    </ul>
</el-form>
</template>

<script>
import backend from "../backend"

export default {
    data() {
        return {
            searchForm: {search: ""},
            avaImpls: [],
            selImpl: null
        }
    },
    methods: {
        async loadImpls() {
            const res = await backend.qryAllImplMods()
            if (typeof res === "string") {
                this.$message.error(`加载可用模块时发生错误：${res}`)
            } else {
                this.avaImpls = res.modSigns
            }
        }
    },
    computed: {
        rowNum() {
            const step = this.avaImpls.length>>2
            let rowIdxs = []
            for (let i = 0; i < this.avaImpls.length; i += step) {
                i = i == 0 ? i : i + 1
                rowIdxs.push(i)
                if (step === 0) {
                    break
                }
            }
            return rowIdxs
        }
    }
}
</script>

<style lang="scss">
.active {
    border-color: #3a8ee6
}
</style>