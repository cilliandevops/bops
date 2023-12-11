<template>
  <div>
    <!-- 用户信息表格 -->
    <a-page-header
    style="border: 3px solid rgb(0, 179, 250)"
    title="Antd表格开发"
    sub-title="cilliandevops"
    @back="() => null"
  />
    <a-table :dataSource="users" :columns="columns" rowKey="id">
      <template #action="{ record }">
        <a-space>
             <!-- 添加用户按钮 -->
          <a-button type="primary" @click="showAddUserModal">添加用户</a-button>
          <a-button type="primary" @click="showEditUserModal(record)">编辑</a-button>
          <a-button type="danger"  @click="deleteUser(record.id)">删除</a-button>
          <a-dropdown>
            <template #overlay>
              <a-menu  @click="handleMenuClick">
                <a-menu-item key="1">编辑</a-menu-item>
                <a-menu-item key="2">删除</a-menu-item>
              </a-menu>
            </template>
            <a-button>
              操作
              <DownOutlined />
            </a-button>
          </a-dropdown>
          <a-button type="primary" loading>加载中</a-button>
        </a-space>
      </template>
    </a-table>


    <!-- 添加用户对话框 -->
    <a-modal
      v-model:visible="addUserModalVisible"
      title="添加用户"
      okText="确认"
      cancelText="取消"
      @ok="addUser"
      @cancel="resetAddUserForm"
    >
      <a-form 
      :model="addUserForm"
      autocomplete="off"
      name="basic"
      
      >
        <a-form-item label="姓名"
        name="name"
        :rules="[{ required: true, message: 'Please input your username!' }]">
          <a-input v-model:value="addUserForm.name" />
        </a-form-item>
        <a-form-item label="年龄"
        
        >
          <a-input-number v-model:value="addUserForm.age" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 编辑用户对话框 -->
    <a-modal
      v-model:visible="editUserModalVisible"
      title="编辑用户"
      okText="确认"
      cancelText="取消"
      @ok="editUser"
    >
      <a-form>
        <a-form-item label="ID" v-if="editUserForm.id">
          <a-input v-model:value="editUserForm.id" disabled />
        </a-form-item>
        <a-form-item label="姓名">
          <a-input v-model:value="editUserForm.name" />
        </a-form-item>
        <a-form-item label="年龄">
          <a-input-number v-model:value="editUserForm.age" />
        </a-form-item>
      </a-form>
    </a-modal>
    <!-- 需要注意的是，这里将id框设置为了禁用状态，以防止用户修改id。同时也需要在编辑用户时将其发送给后端 -->
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from 'vue'
import axios from 'axios'
import { DownOutlined } from '@ant-design/icons-vue';
import type { MenuProps } from 'ant-design-vue';


interface FormState {
  name: string;
  age: string;

}

export default defineComponent({
  components: {
    DownOutlined,
  },

  setup() {

    const formState = reactive<FormState>({
      name: '',
      age: '',
  
    });

    const handleMenuClick: MenuProps['onClick'] = e => {
      console.log('click', e);
    };

    // 用户数据
    const users = ref([])

    // 获取用户数据
    const fetchUsers = async () => {
      const response = await axios.get('http://localhost:8089/users')
      users.value = response.data
    }

    // 表格列配置
    const columns = [
      { title: 'ID', dataIndex: 'id', key: 'id' },
      { title: '姓名', dataIndex: 'name', key: 'name' },
      { title: '年龄', dataIndex: 'age', key: 'age' },
      {
        title: '操作',
        key: 'action',
        slots: { customRender: 'action' },
      },
    ]
    

    // 添加用户对话框
    const addUserModalVisible = ref(false)
    const addUserForm = reactive({ name: '', age: 0 })

    // 显示添加用户对话框
    const showAddUserModal = () => {
      addUserModalVisible.value = true
    }

    // 添加用户
    const addUser = async () => { 
      try {
        await axios.post('http://localhost:8089/users', addUserForm)
      } catch (err) {
        console.error(err)
      }
      addUserModalVisible.value = false
      fetchUsers()
      resetAddUserForm()
    }

    // 重置添加用户表单
    const resetAddUserForm = () => {
      addUserForm.name = ''
      addUserForm.age = 0
    }

    // 编辑用户对话框
    const editUserModalVisible = ref(false)
    const editUserForm = reactive({ id: 0, name: '', age: 0 })

    // 显示编辑用户对话框
    const showEditUserModal = (user) => {
      editUserForm.id = user.id
      editUserForm.name = user.name
      editUserForm.age = user.age
      editUserModalVisible.value = true
    }

    // 编辑用户
    const editUser = async () => {
      await axios.put(`http://localhost:8089/users/${editUserForm.id}`, editUserForm)
      editUserModalVisible.value = false
      fetchUsers()
    }

    // 删除用户
    const deleteUser = async (id) => {
      await axios.delete(`http://localhost:8089/users/${id}`)
      fetchUsers()
    }

    // 初始化
    fetchUsers()

    return {
      
      formState,
      handleMenuClick,
      users,
      columns,
      addUserModalVisible,
      addUserForm,
      showAddUserModal,
      addUser,
      resetAddUserForm,
      editUserModalVisible,
      editUserForm,
      showEditUserModal,
      editUser,
      deleteUser,
    }
  },
})
</script>