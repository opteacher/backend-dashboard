<template>
<el-row class="toolbar">
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-left" size="mini"/>
    </el-col>
    <el-col class="p-10" :span="22">
        <el-button-group class="p-0">
            <el-button class="p-7" type="primary" icon="el-icon-menu" size="mini" @click="showSelItfDlg = true"/>
        </el-button-group>
    </el-col>
    <el-col class="p-10" :span="1">
        <el-button class="p-7" plain icon="el-icon-arrow-right" size="mini"/>
    </el-col>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="选择接口" :visible.sync="showSelItfDlg" :modal-append-to-body="false" width="40vw">
        <sel-interface ref="sel-itf-form"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showSelItfDlg = false">取 消</el-button>
            <el-button type="primary">确 定</el-button>
        </div>
    </el-dialog>
</el-row>
</template>

<script>
import apisBkd from "../async/api"
import selInterface from "../forms/selInterface"

export default {
    components: {
        "sel-interface": selInterface
    },
    data() {
        return {
            showSelItfDlg: false
        }
    },
    async created() {
        console.log(await apisBkd.qry())
    }
}
</script>

