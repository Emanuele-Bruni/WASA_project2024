import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js'
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import UserItem from './components/UserItem.vue'
import SuccessMsg from './components/SuccessMsg.vue'
import ConversationItem from './components/ConversationItem.vue'
import MessageItem from './components/MessageItem.vue'
import ChatInfoItem from './components/ChatInfoItem.vue'
import ConfirmModal from './components/ConfirmModal.vue'
import MessageModal from './components/MessageModal.vue'
import ForwardModal from './components/ForwardModal.vue'
import AddToGroupModal from './components/AddToGroupModal.vue'
import EmojiPicker from './components/EmojiPicker.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("UserItem",UserItem);
app.component("SuccessMsg", SuccessMsg);
app.component("ConversationItem",ConversationItem);
app.component("MessageItem", MessageItem);
app.component("ChatInfoItem", ChatInfoItem);
app.component("ConfirmModal", ConfirmModal);
app.component("MessageModal", MessageModal);
app.component("ForwardModal", ForwardModal);
app.component("AddToGroupModal", AddToGroupModal);
app.component("EmojiPicker", EmojiPicker);
app.use(router)
app.mount('#app')
