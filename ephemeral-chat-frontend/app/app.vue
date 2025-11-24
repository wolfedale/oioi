<template>
  <div class="h-screen w-screen bg-gradient-to-br from-gray-900 via-gray-800 to-indigo-900 relative overflow-hidden" style="min-height: 100vh; background: linear-gradient(135deg, #1f2937 0%, #374151 50%, #1e40af 100%);">
    <NuxtRouteAnnouncer />

    <!-- CONTENT ABSOLUTELY CENTERED WITH INLINE STYLES -->
    <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); width: 100%; padding: 0 1.5rem;">
      <!-- Login Container (narrow) -->
      <div v-if="!currentUser" class="w-full max-w-md mx-auto">
        <!-- Account Creation Section -->
        <div class="backdrop-blur-xl rounded-xl p-8 mb-6 text-center animate-fade-in shadow-2xl bg-gray-900/90 border border-gray-700/30">
          <div class="mb-6">
            <h1 class="text-3xl font-bold text-gradient mb-2">oi.oi</h1>
            <p class="text-sm text-gray-400">Messages that exist only in the moment</p>
          </div>

          <!-- Step 1: Create Account or Login -->
          <div v-if="!hasAccount" class="space-y-4">
            <!-- Toggle between Create/Login -->
            <div class="flex bg-gray-800/60 rounded-lg p-1 mb-4">
              <button
                @click="isLoginMode = false"
                :class="[
                  'flex-1 py-2 px-4 rounded-md transition-all duration-200 text-sm font-medium',
                  !isLoginMode ? 'bg-blue-600 text-white shadow-lg' : 'text-gray-300 hover:text-white hover:bg-gray-700/40'
                ]"
              >
                Create Account
              </button>
              <button
                @click="isLoginMode = true"
                :class="[
                  'flex-1 py-2 px-4 rounded-md transition-all duration-200 text-sm font-medium',
                  isLoginMode ? 'bg-blue-600 text-white shadow-lg' : 'text-gray-300 hover:text-white hover:bg-gray-700/40'
                ]"
              >
                Login
              </button>
            </div>

            <!-- Create Account Mode -->
            <div v-if="!isLoginMode">
              <div class="mb-4">
                <h3 class="text-lg font-medium mb-2 text-white">Create Account</h3>
                <p class="text-xs text-gray-400">Write a sentence to generate your unique ID</p>
              </div>
              <textarea
                v-model="sentence"
                placeholder="Write a unique sentence (minimum 10 characters)..."
                class="w-full px-4 py-3 rounded-lg bg-gray-800/70 border border-gray-600/50 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:bg-gray-800 transition-all duration-200 font-medium text-center resize-none h-20"
                @keyup.ctrl.enter="createAccount"
              ></textarea>
              <button
                @click="createAccount"
                :disabled="!sentence.trim() || sentence.trim().length < 10"
                class="w-full mt-4 px-6 py-3 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-medium rounded-lg transition-all duration-200 shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:shadow-lg"
              >
                Generate ID
              </button>
              <p class="text-gray-500 text-xs mt-2">Tip: Remember your sentence - it's your only way back!</p>
            </div>

            <!-- Login Mode -->
            <div v-else>
              <div class="mb-4">
                <h3 class="text-lg font-medium text-white mb-2">Login</h3>
                <p class="text-gray-400 text-xs">Enter your ID and sentence to login</p>
              </div>
              <input
                v-model="loginUserID"
                type="text"
                placeholder="Your 8-character ID..."
                class="w-full px-4 py-3 rounded-lg bg-gray-800/70 border border-gray-600/50 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:bg-gray-800 transition-all duration-200 font-medium text-center mb-3 font-mono"
                maxlength="8"
              />
              <textarea
                v-model="sentence"
                placeholder="Your original sentence..."
                class="w-full px-4 py-3 rounded-lg bg-gray-800/70 border border-gray-600/50 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:bg-gray-800 transition-all duration-200 font-medium text-center resize-none h-20"
                @keyup.ctrl.enter="loginAccount"
              ></textarea>
              <button
                @click="loginAccount"
                :disabled="!loginUserID.trim() || !sentence.trim() || sentence.trim().length < 10"
                class="w-full mt-4 px-6 py-3 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-medium rounded-lg transition-all duration-200 shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:shadow-lg"
              >
                Login
              </button>
              <p class="text-gray-500 text-xs mt-2">Both ID and sentence must match exactly</p>
            </div>
          </div>

          <!-- Step 2: Show Generated ID -->
          <div v-else class="space-y-4">
            <div class="mb-4">
              <h3 class="text-lg font-medium text-white mb-2">Your Account ID</h3>
              <p class="text-gray-400 text-xs">This is your unique chat identifier</p>
            </div>
            <div class="bg-gray-800/80 border border-gray-600 rounded-lg p-4 mb-4 transition-colors duration-300">
              <p class="text-2xl font-mono font-bold text-blue-400 tracking-widest">{{ generatedUserID }}</p>
            </div>
            <button
              @click="useAccount"
              class="w-full px-6 py-3 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-medium rounded-lg transition-all duration-200 shadow-lg hover:shadow-xl"
            >
              Start Chatting
            </button>
            <p class="text-gray-500 text-xs">Save this ID - you'll need your sentence to regenerate it</p>
          </div>
        </div>
      </div>

      <!-- Chat Container (compact) -->
      <div v-if="currentUser" class="w-full max-w-4xl mx-auto">
        <!-- Clean Chat Interface -->
        <div v-if="currentUser" class="backdrop-blur-xl rounded-2xl flex flex-col h-[75vh] animate-slide-up overflow-hidden shadow-2xl bg-gray-900/95 border border-gray-700/30">
          <!-- Clean Header -->
          <div class="flex items-center justify-between px-6 py-4 bg-gray-800/60 border-b border-gray-700/40">
            <div class="flex items-center space-x-3">
              <div class="relative">
                <div class="w-7 h-7 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-bold text-xs shadow-md">
                  {{ currentUser.slice(0, 2).toUpperCase() }}
                </div>
                <div class="absolute -bottom-0.5 -right-0.5 w-2.5 h-2.5 rounded-full bg-green-500 border-2 border-white dark:border-gray-700 animate-pulse"></div>
              </div>
              <div>
                <h2 class="text-sm font-semibold text-white">{{ currentUser }}</h2>
                <p class="text-xs text-gray-400">Online • oi-oi</p>
              </div>
            </div>
            <div class="flex items-center space-x-3">
              <button
                @click="showNewChatModal = true"
                class="flex items-center px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-all duration-200 font-medium shadow-lg hover:shadow-xl text-sm"
                title="New Conversation"
              >
                <Plus class="w-4 h-4" />
              </button>
              <button
                @click="leaveChat"
                class="p-2 text-gray-400 hover:text-red-400 hover:bg-red-900/20 rounded-lg transition-all duration-200"
                title="Sign Out"
              >
                <LogOut class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Conversation Tabs -->
          <div class="flex items-center space-x-1 px-6 py-3 overflow-x-auto bg-gray-800/40 border-b border-gray-700/40">
            <div
              v-for="[conversationId, conversation] in conversations"
              :key="conversationId"
              @click="switchToConversation(conversationId)"
              :class="[
                'flex items-center space-x-2 px-2.5 py-1 cursor-pointer transition-all duration-200 whitespace-nowrap group rounded-lg',
                activeConversationId === conversationId
                  ? 'bg-gray-600 text-blue-400 shadow-sm border border-blue-400/30'
                  : 'text-slate-300 hover:text-slate-100 hover:bg-gray-600/50'
              ]"
            >
              <div class="relative">
                <div class="w-5 h-5 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white text-xs font-bold">
                  {{ conversationId.slice(0, 2).toUpperCase() }}
                </div>
                <!-- New message indicator -->
                <div
                  v-if="hasUnreadMessages(conversationId)"
                  class="absolute -top-0.5 -right-0.5 w-2 h-2 bg-red-500 rounded-full border border-white dark:border-gray-700 animate-pulse"
                  title="New messages"
                ></div>
              </div>
              <div class="flex flex-col min-w-0">
                <span class="font-medium text-xs truncate">{{ conversationId }}</span>
              </div>

              <button
                @click.stop="removeConversation(conversationId)"
                class="opacity-0 group-hover:opacity-100 text-slate-500 hover:text-red-400 transition-all duration-200 p-0.5 rounded-md hover:bg-red-900/30"
                title="Close conversation"
              >
                <X class="w-2.5 h-2.5" />
              </button>
            </div>
          </div>

          <!-- Clean Chat Area -->
          <div class="flex-1 flex flex-col min-h-0 bg-gray-900/50">
            <div v-if="!activeConversationId" class="flex-1 flex items-center justify-center">
              <div class="text-center p-8">
                <div class="w-16 h-16 bg-gradient-to-br from-blue-600/20 to-purple-600/20 rounded-full flex items-center justify-center mx-auto mb-4">
                  <MessageCircle class="w-8 h-8 text-blue-400" />
                </div>
                <h2 class="text-lg font-semibold mb-2 text-white">Welcome to oi-oi</h2>
                <p class="text-sm text-gray-400">Select a conversation or use the + button to start a new one</p>
              </div>
            </div>

            <div v-else class="flex-1 flex flex-col min-h-0">
              <!-- Messages Container -->
              <div
                ref="messagesContainer"
                class="flex-1 overflow-y-auto scroll-smooth p-4 space-y-3"
              >
                <!-- Empty State -->
                <div v-if="messages.length === 0" class="text-center py-6">
                  <div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-900/50 to-purple-900/50 flex items-center justify-center mx-auto mb-3">
                    <MessageCircle class="w-5 h-5 text-blue-400" />
                  </div>
                  <h3 class="text-sm font-semibold text-slate-200 mb-1">Start the conversation</h3>
                  <p class="text-xs text-slate-400">Send a message to begin your ephemeral chat</p>
                </div>

                <!-- Messages -->
                <div
                  v-for="message in messages"
                  :key="message.id"
                  class="animate-fade-in message-container"
                  :class="{ 'message-fading': isMessageExpiring(message) }"
                >
                  <div
                    v-if="message.from === currentUser"
                    class="flex justify-end mb-2"
                  >
                    <div class="max-w-xs lg:max-w-md">
                      <div class="bg-blue-500 text-white px-3 py-2 rounded-xl rounded-br-md shadow-md transition-colors duration-300">
                        <p class="text-xs leading-relaxed">{{ message.content }}</p>
                      </div>
                      <p class="text-xs text-slate-400 mt-0.5 text-right">
                        {{ formatTime(message.timestamp) }}
                      </p>
                    </div>
                  </div>

                  <div v-else class="flex justify-start mb-2">
                    <div class="flex items-start space-x-2 max-w-xs lg:max-w-md">
                      <div class="w-5 h-5 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white text-xs font-bold flex-shrink-0">
                        {{ message.from.slice(0, 2).toUpperCase() }}
                      </div>
                      <div class="flex-1">
                        <div class="bg-gray-600 border border-gray-500 px-3 py-2 rounded-xl rounded-bl-md shadow-sm transition-colors duration-300">
                          <p class="text-xs font-medium text-slate-300 mb-0.5">{{ message.from }}</p>
                          <p class="text-xs text-slate-100 leading-relaxed">{{ message.content }}</p>
                        </div>
                        <p class="text-xs text-slate-500 mt-0.5">
                          {{ formatTime(message.timestamp) }}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Clean Message Input -->
              <div class="p-6 bg-gray-800/70 border-t border-gray-700/50">
                <div class="flex items-center space-x-3">
                  <div class="flex-1 relative">
                    <input
                      v-model="newMessage"
                      type="text"
                      :placeholder="isSending ? 'Sending...' : 'Type your message...'"
                      class="w-full px-4 py-2.5 rounded-xl bg-gray-700 border border-gray-600 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent focus:bg-gray-600 transition-all duration-200 text-xs pr-12"
                      :class="{ 'opacity-50 cursor-not-allowed': isSending }"
                      @keyup.enter="sendMessage"
                      :disabled="!activeConversationId || isSending"
                    />
                    <button
                      @click="sendMessage"
                      :disabled="!newMessage.trim() || isSending"
                      class="absolute right-2 top-1/2 transform -translate-y-1/2 w-7 h-7 rounded-lg transition-all duration-200 flex items-center justify-center"
                      :class="newMessage.trim() && !isSending
                        ? 'bg-blue-600 hover:bg-blue-700 text-white shadow-md hover:shadow-lg'
                        : isSending
                        ? 'bg-blue-400 text-white cursor-not-allowed'
                        : 'bg-gray-600 text-slate-500 cursor-not-allowed'"
                    >
                      <Send class="w-3.5 h-3.5" :class="{ 'animate-pulse': isSending }" />
                    </button>
                  </div>
                </div>
                <div class="flex items-center justify-between mt-2 px-1">
                  <p class="text-xs text-slate-500">Messages disappear after 10 seconds</p>
                  <div class="flex items-center space-x-1 text-xs text-slate-500">
                    <div class="w-2 h-2 rounded-full bg-green-500"></div>
                    <span>Connected</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="text-center mt-4 text-xs text-slate-400">
          <p>End-to-end ephemeral messaging • Zero data retention</p>
        </div>
      </div>
    </div>

    <!-- Copyright Footer - Fixed at bottom with inline styles -->
    <div style="position: fixed; bottom: 0; left: 0; right: 0; z-index: 9999; background: rgba(31, 41, 55, 0.9); padding: 12px 0;">
      <div style="text-align: center; font-size: 0.75rem; color: #9ca3af;">
        <p style="margin: 0;">© 2024 REDGRIP.io. All rights reserved.</p>
        <p style="margin: 4px 0 0 0;">Powered by REDGRIP Tier</p>
      </div>
    </div>

    <!-- New Chat Modal -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="opacity-0 transform scale-95"
      enter-to-class="opacity-100 transform scale-100"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="opacity-100 transform scale-100"
      leave-to-class="opacity-0 transform scale-95"
    >
      <div
        v-if="showNewChatModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black bg-opacity-50"
          @click="showNewChatModal = false"
        ></div>

        <!-- Modal -->
        <div class="relative bg-gray-800 rounded-2xl p-8 w-full max-w-md shadow-2xl border border-gray-600 transition-colors duration-300">
          <div class="flex items-center justify-between mb-6">
            <div>
              <h3 class="text-xl font-semibold text-slate-200">New Conversation</h3>
              <p class="text-sm text-slate-400 mt-1">Connect with another user</p>
            </div>
            <button
              @click="showNewChatModal = false"
              class="text-slate-500 hover:text-slate-300 transition-colors p-2 rounded-lg hover:bg-gray-700"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <div class="mb-6">
            <label class="block text-sm font-semibold text-slate-300 mb-3">User ID</label>
            <input
              v-model="newChatUserId"
              type="text"
              placeholder="Enter 8-character user ID..."
              class="w-full px-4 py-3 rounded-xl bg-gray-700 border border-gray-600 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent focus:bg-gray-600 transition-all duration-200 font-mono text-center"
              maxlength="8"
              @keyup.enter="startNewConversation"
            />
            <p class="text-xs text-slate-400 mt-2 text-center">Only users with exact matching IDs can chat</p>
          </div>

          <div class="flex space-x-3">
            <button
              @click="showNewChatModal = false"
              class="flex-1 px-6 py-3 text-slate-300 bg-gray-700 hover:bg-gray-600 rounded-xl transition-all duration-200 font-medium"
            >
              Cancel
            </button>
            <button
              @click="startNewConversation"
              :disabled="!newChatUserId.trim()"
              class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-xl transition-all duration-200 font-medium shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:shadow-lg"
            >
              Start Chat
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Notification Toast -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="opacity-0 transform translate-y-2"
      enter-to-class="opacity-100 transform translate-y-0"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="opacity-100 transform translate-y-0"
      leave-to-class="opacity-0 transform translate-y-2"
    >
      <div
        v-if="notification.visible"
        class="fixed top-4 left-1/2 transform -translate-x-1/2 z-50 max-w-md w-full mx-4"
      >
        <div
          :class="[
            'notification-toast',
            notification.type === 'warning' ? 'notification-warning' : '',
            notification.type === 'error' ? 'notification-error' : '',
            notification.type === 'success' ? 'notification-success' : '',
            notification.type === 'info' ? 'notification-info' : ''
          ]"
        >
          <div class="flex items-center space-x-3">
            <!-- Icon based on type -->
            <div class="flex-shrink-0">
              <div v-if="notification.type === 'warning'" class="w-5 h-5 text-amber-500">
                ⚠️
              </div>
              <div v-else-if="notification.type === 'error'" class="w-5 h-5 text-red-500">
                ❌
              </div>
              <div v-else-if="notification.type === 'success'" class="w-5 h-5 text-green-500">
                ✅
              </div>
              <div v-else class="w-5 h-5 text-blue-500">
                ℹ️
              </div>
            </div>

            <!-- Message -->
            <div class="flex-1">
              <p class="text-sm font-medium text-white">{{ notification.message }}</p>
            </div>

            <!-- Close button -->
            <button
              @click="notification.visible = false"
              class="flex-shrink-0 text-gray-400 hover:text-gray-200 transition-colors"
            >
              <span class="sr-only">Close</span>
              ✕
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { MessageCircle, Send, LogOut, Plus, X, Sun, Moon } from 'lucide-vue-next'

// Types
interface Message {
  id: string
  from: string
  to: string
  content: string
  timestamp: string
  expires_at: string
  readAt?: string // When the message was read
}

interface Conversation {
  id: string // recipient user ID
  lastMessage?: Message
  lastActivity: string
  messages: Message[]
  messageTimeouts: Map<string, NodeJS.Timeout>
}

// Reactive state
const currentUser = ref<string>('')
const newMessage = ref<string>('')
const messagesContainer = ref<HTMLElement>()
const isSending = ref<boolean>(false)

// Conversation management
const conversations = ref<Map<string, Conversation>>(new Map())
const activeConversationId = ref<string>('')
const showNewChatModal = ref<boolean>(false)
const newChatUserId = ref<string>('')


// Computed for current conversation
const activeConversation = computed(() =>
  activeConversationId.value ? conversations.value.get(activeConversationId.value) : null
)

const messages = computed(() => activeConversation.value?.messages || [])


// Check if conversation has unread messages
const hasUnreadMessages = (conversationId: string) => {
  const conversation = conversations.value.get(conversationId)
  if (!conversation || !conversation.lastMessage) return false

  // If this conversation is not active and has messages, consider them unread
  return activeConversationId.value !== conversationId && conversation.messages.length > 0
}

// Account creation state
const hasAccount = ref<boolean>(false)
const sentence = ref<string>('')
const generatedUserID = ref<string>('')
const isLoginMode = ref<boolean>(false)
const loginUserID = ref<string>('')

// Notification system
const notification = ref<{
  message: string
  type: 'info' | 'warning' | 'error' | 'success'
  visible: boolean
}>({
  message: '',
  type: 'info',
  visible: false
})


// Runtime config
const config = useRuntimeConfig()

let presenceInterval: NodeJS.Timeout | null = null
let uiUpdateInterval: NodeJS.Timeout | null = null

// Reactive timestamp for UI updates
const currentTime = ref(Date.now())

// Conversation management methods
const getOrCreateConversation = (recipientId: string): Conversation => {
  if (!conversations.value.has(recipientId)) {
    const conversation: Conversation = {
      id: recipientId,
      lastActivity: new Date().toISOString(),
      messages: [],
      messageTimeouts: new Map()
    }
    conversations.value.set(recipientId, conversation)
  }
  return conversations.value.get(recipientId)!
}

const switchToConversation = (recipientId: string) => {
  activeConversationId.value = recipientId
  // Scroll to bottom when switching conversations
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const startNewConversation = () => {
  if (!newChatUserId.value.trim()) return

  const recipientId = newChatUserId.value.trim()
  getOrCreateConversation(recipientId)
  switchToConversation(recipientId)

  // Clear modal
  showNewChatModal.value = false
  newChatUserId.value = ''
}

const removeConversation = (conversationId: string) => {
  const conversation = conversations.value.get(conversationId)
  if (!conversation) return

  // Clear all timeouts for this conversation
  conversation.messageTimeouts.forEach(timeout => clearTimeout(timeout))
  conversation.messageTimeouts.clear()

  // Remove conversation
  conversations.value.delete(conversationId)

  // If this was the active conversation, switch to another or clear
  if (activeConversationId.value === conversationId) {
    const remainingConversations = Array.from(conversations.value.keys())
    if (remainingConversations.length > 0) {
      switchToConversation(remainingConversations[0])
    } else {
      activeConversationId.value = ''
    }
  }
}

// Methods
const showNotification = (message: string, type: 'info' | 'warning' | 'error' | 'success' = 'info') => {
  notification.value = {
    message,
    type,
    visible: true
  }

  // Auto-hide after 4 seconds
  setTimeout(() => {
    notification.value.visible = false
  }, 4000)
}

const createAccount = async () => {
  if (!sentence.value.trim()) return

  try {
    const response = await fetch(`${config.public.apiBase}/account/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ sentence: sentence.value.trim() })
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || 'Failed to create account')
    }

    generatedUserID.value = data.user_id
    hasAccount.value = true

    // Clear sensitive sentence from memory
    sentence.value = ''

  } catch (error: any) {
    console.error('Account creation error:', error)
    showNotification(error.message || 'Failed to create account. Please try again.', 'error')
  }
}

const loginAccount = async () => {
  if (!loginUserID.value.trim() || !sentence.value.trim()) return

  try {
    const response = await fetch(`${config.public.apiBase}/account/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: loginUserID.value.trim(),
        sentence: sentence.value.trim()
      })
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || 'Login failed')
    }

    generatedUserID.value = data.user_id
    hasAccount.value = true

    // Clear sensitive data from memory
    sentence.value = ''
    loginUserID.value = ''

  } catch (error: any) {
    console.error('Login error:', error)
    showNotification(error.message || 'Login failed. Please check your credentials.', 'error')
  }
}

const useAccount = () => {
  currentUser.value = generatedUserID.value

  // Start presence pinging
  startPresence()

  // Connect to WebSocket for real-time messaging
  connectWebSocket()

  // Start UI update timer for fade effects
  uiUpdateInterval = setInterval(() => {
    currentTime.value = Date.now()
  }, 1000) // Update every second
}

const scheduleMessageRemoval = (messageId: string, conversationId: string) => {
  const conversation = conversations.value.get(conversationId)
  if (!conversation) return

  // Clear existing timeout if any
  const existingTimeout = conversation.messageTimeouts.get(messageId)
  if (existingTimeout) {
    clearTimeout(existingTimeout)
  }

  // Set new timeout for 10 seconds
  const timeout = setTimeout(() => {
    console.log(`Removing message ${messageId} after 10 seconds`)
    conversation.messages = conversation.messages.filter(msg => msg.id !== messageId)
    conversation.messageTimeouts.delete(messageId)
  }, 10000) // 10 seconds

  conversation.messageTimeouts.set(messageId, timeout)
}


const leaveChat = () => {
  // Clear intervals
  if (presenceInterval) clearInterval(presenceInterval)
  if (uiUpdateInterval) clearInterval(uiUpdateInterval)

  // Close WebSocket connection
  if (websocket) {
    websocket.close()
    websocket = null
  }

  // Clear all conversation timeouts
  conversations.value.forEach(conversation => {
    conversation.messageTimeouts.forEach(timeout => clearTimeout(timeout))
    conversation.messageTimeouts.clear()
  })

  // Reset state
  currentUser.value = ''
  conversations.value.clear()
  activeConversationId.value = ''
  hasAccount.value = false
  sentence.value = ''
  generatedUserID.value = ''
  isLoginMode.value = false
  loginUserID.value = ''
}


const startPresence = async () => {
  const pingPresence = async () => {
    try {
      console.log(`Pinging presence for user: ${currentUser.value}`)
      const response = await fetch(`${config.public.apiBase}/presence/ping`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: currentUser.value })
      })
      const data = await response.json()
      console.log('Presence ping successful:', data)
    } catch (error) {
      console.error('Presence ping failed:', error)
    }
  }

  // Initial ping and fetch online users
  await pingPresence()

  // Set up interval
  presenceInterval = setInterval(pingPresence, 15000) // Every 15 seconds
}

// WebSocket connection for real-time messaging
let websocket: WebSocket | null = null

const connectWebSocket = () => {
  if (!currentUser.value) return

  const wsUrl = `ws://localhost:3001/ws?user_id=${currentUser.value}`
  console.log(`🔌 Connecting to WebSocket: ${wsUrl}`)

  websocket = new WebSocket(wsUrl)

  websocket.onopen = () => {
    console.log('✅ WebSocket connected')
  }

  websocket.onmessage = (event) => {
    try {
      const wsMessage = JSON.parse(event.data)
      console.log('📨 WebSocket message received:', wsMessage)

      switch (wsMessage.type) {
        case 'new_message':
          handleIncomingMessage(wsMessage)
          break
        case 'message_sent':
          console.log('✅ Message sent confirmation:', wsMessage.message_id)
          break
        case 'pong':
          console.log('🏓 WebSocket pong received')
          break
        case 'error':
          console.error('❌ WebSocket error:', wsMessage.error)
          showNotification(wsMessage.error, 'error')
          break
      }
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error)
    }
  }

  websocket.onclose = (event) => {
    console.log('🔌 WebSocket disconnected:', event.code, event.reason)

    // Attempt to reconnect after 3 seconds
    if (currentUser.value) {
      setTimeout(() => {
        console.log('🔄 Attempting WebSocket reconnect...')
        connectWebSocket()
      }, 3000)
    }
  }

  websocket.onerror = (error) => {
    console.error('❌ WebSocket error:', error)
  }

  // Send ping every 30 seconds to keep connection alive
  const pingInterval = setInterval(() => {
    if (websocket?.readyState === WebSocket.OPEN) {
      websocket.send(JSON.stringify({ type: 'ping' }))
    } else {
      clearInterval(pingInterval)
    }
  }, 30000)
}

const handleIncomingMessage = (wsMessage: any) => {
  console.log(`📨 New message: ${wsMessage.message_id} from ${wsMessage.from}`)

  // Add readAt timestamp
  const messageWithReadTime = {
    id: wsMessage.message_id,
    from: wsMessage.from,
    to: wsMessage.to,
    content: wsMessage.content,
    timestamp: wsMessage.timestamp,
    expires_at: wsMessage.expires_at,
    readAt: new Date().toISOString()
  }

  // Determine which conversation this message belongs to
  const conversationId = messageWithReadTime.from === currentUser.value
    ? messageWithReadTime.to
    : messageWithReadTime.from

  // Get or create conversation
  const conversation = getOrCreateConversation(conversationId)

  // Add message to conversation
  conversation.messages.push(messageWithReadTime)
  conversation.lastMessage = messageWithReadTime
  conversation.lastActivity = new Date().toISOString()

  // Schedule removal after 10 seconds
  scheduleMessageRemoval(messageWithReadTime.id, conversationId)

  console.log(`💬 Message added to conversation ${conversationId}: "${messageWithReadTime.content}"`)

  // Auto-scroll to bottom if this is the active conversation
  if (activeConversationId.value === conversationId) {
    nextTick(() => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    })
  }
}


// Robust API call with retry logic
const apiCall = async (url: string, options: RequestInit, retries = 2): Promise<Response> => {
  for (let i = 0; i <= retries; i++) {
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options.headers
        }
      })
      return response
    } catch (error) {
      if (i === retries) throw error
      console.log(`API call failed, retrying... (${i + 1}/${retries})`)
      await new Promise(resolve => setTimeout(resolve, 500 * (i + 1))) // Progressive delay
    }
  }
  throw new Error('All retries failed')
}

const sendMessage = async () => {
  if (!newMessage.value.trim() || !activeConversationId.value || isSending.value) return

  // Check WebSocket connection
  if (!websocket || websocket.readyState !== WebSocket.OPEN) {
    showNotification('WebSocket connection lost. Reconnecting...', 'warning')
    connectWebSocket()
    return
  }

  // Prevent rapid-fire sends
  isSending.value = true

  const targetUser = activeConversationId.value
  const messageContent = newMessage.value.trim()

  try {
    console.log(`📤 Sending WebSocket message to ${targetUser}: "${messageContent}"`)

    // Send message via WebSocket
    const wsMessage = {
      type: 'send_message',
      to: targetUser,
      content: messageContent
    }

    websocket.send(JSON.stringify(wsMessage))

    // Clear input immediately since WebSocket is synchronous
    newMessage.value = ''

    // Add sent message immediately to conversation for instant feedback
    const sentMessage = {
      id: `sent-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      from: currentUser.value,
      to: targetUser,
      content: messageContent,
      timestamp: new Date().toISOString(),
      expires_at: new Date(Date.now() + 60000).toISOString(),
      readAt: new Date().toISOString()
    }

    // Get the conversation and add the message
    const conversation = getOrCreateConversation(targetUser)
    conversation.messages.push(sentMessage)
    conversation.lastMessage = sentMessage
    conversation.lastActivity = new Date().toISOString()

    scheduleMessageRemoval(sentMessage.id, targetUser)
    console.log(`💬 Sent message displayed locally: "${messageContent}"`)

    // Auto-scroll to bottom
    nextTick(() => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    })

  } catch (error: any) {
    console.error('❌ Send message error:', error)
    showNotification('Failed to send message. Please try again.', 'error')

    // Restore message in input if sending failed
    newMessage.value = messageContent
  } finally {
    // Always re-enable sending after attempt
    isSending.value = false
  }
}

const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const isMessageExpiring = (message: Message) => {
  if (!message.readAt) return false

  const readTime = new Date(message.readAt).getTime()
  const now = currentTime.value // Use reactive currentTime
  const elapsedSeconds = (now - readTime) / 1000

  // Start fading after 7 seconds (3 seconds before disappearing)
  return elapsedSeconds > 7
}

// Cleanup on unmount
onUnmounted(() => {
  if (presenceInterval) clearInterval(presenceInterval)
  if (uiUpdateInterval) clearInterval(uiUpdateInterval)

  // Close WebSocket connection
  if (websocket) {
    websocket.close()
    websocket = null
  }

  // Clear all conversation timeouts
  conversations.value.forEach(conversation => {
    conversation.messageTimeouts.forEach(timeout => clearTimeout(timeout))
    conversation.messageTimeouts.clear()
  })
})
</script>

<style>
@import 'tailwindcss/base';
@import 'tailwindcss/components';
@import 'tailwindcss/utilities';

@layer base {
  html {
    font-family: 'Inter', system-ui, -apple-system, sans-serif;
    font-feature-settings: 'cv02', 'cv03', 'cv04', 'cv11';
  }

  body {
    @apply text-slate-100 antialiased;
    min-height: 100vh;
    background: linear-gradient(135deg, #1f2937 0%, #374151 50%, #1e40af 100%);
  }
}

@layer components {
  .glass {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.2);
  }

  .glass-dark {
    background: rgba(0, 0, 0, 0.3);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .btn-primary {
    @apply bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700;
    @apply text-white font-medium px-6 py-3 rounded-lg;
    @apply transform transition-all duration-200 ease-out;
    @apply hover:scale-105 hover:shadow-xl;
    @apply focus:outline-none focus:ring-4 focus:ring-blue-500/25;
  }

  .btn-primary-dark {
    @apply bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700;
    @apply text-white font-medium px-6 py-3 rounded-lg;
    @apply transform transition-all duration-200 ease-out;
    @apply hover:scale-105 hover:shadow-xl;
    @apply focus:outline-none focus:ring-4 focus:ring-blue-500/25;
  }

  .input-modern {
    @apply w-full px-4 py-3 rounded-lg;
    @apply bg-white/50 backdrop-blur-sm;
    @apply border border-slate-200/50;
    @apply text-slate-900 placeholder-slate-500;
    @apply focus:outline-none focus:ring-4 focus:ring-blue-500/25;
    @apply focus:border-blue-500/50 focus:bg-white/80;
    @apply transition-all duration-200;
  }

  .input-modern-dark {
    @apply w-full px-4 py-3 rounded-lg;
    @apply bg-gray-800/50 backdrop-blur-sm;
    @apply border border-gray-600/50;
    @apply text-white placeholder-gray-400;
    @apply focus:outline-none focus:ring-4 focus:ring-blue-500/25;
    @apply focus:border-blue-500/50 focus:bg-gray-800/80;
    @apply transition-all duration-200;
  }

  .chat-bubble-sent {
    @apply bg-gradient-to-r from-blue-500 to-purple-600;
    @apply text-white px-4 py-3 rounded-lg rounded-br-md;
    @apply max-w-xs ml-auto;
    @apply shadow-lg shadow-blue-500/25;
  }

  .chat-bubble-sent-dark {
    @apply bg-gradient-to-r from-blue-600 to-purple-600;
    @apply text-white px-4 py-3 rounded-lg rounded-br-md;
    @apply max-w-xs ml-auto;
    @apply shadow-lg shadow-blue-500/25;
  }

  .chat-bubble-received {
    @apply bg-white/80 backdrop-blur-sm;
    @apply text-slate-900 px-4 py-3 rounded-lg rounded-bl-md;
    @apply max-w-xs mr-auto;
    @apply border border-slate-200/50;
    @apply shadow-lg shadow-slate-500/10;
  }

  .chat-bubble-received-dark {
    @apply bg-gray-800/80 backdrop-blur-sm;
    @apply text-white px-4 py-3 rounded-lg rounded-bl-md;
    @apply max-w-xs mr-auto;
    @apply border border-gray-600/50;
    @apply shadow-lg shadow-gray-900/25;
  }

  .scroll-smooth {
    scrollbar-width: thin;
    scrollbar-color: rgb(203 213 225) transparent;
  }

  .scroll-smooth::-webkit-scrollbar {
    width: 6px;
  }

  .scroll-smooth::-webkit-scrollbar-track {
    background: transparent;
  }

  .scroll-smooth::-webkit-scrollbar-thumb {
    background: rgb(203 213 225);
    border-radius: 3px;
  }

  .scroll-smooth::-webkit-scrollbar-thumb:hover {
    background: rgb(148 163 184);
  }
}

@layer utilities {
  .animate-fade-in {
    animation: fadeIn 0.5s ease-out forwards;
  }

  .animate-slide-up {
    animation: slideUp 0.4s ease-out forwards;
  }

  .animate-pulse-slow {
    animation: pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite;
  }

  /* Message expiring animation */
  .message-container {
    transition: opacity 3s ease-out, transform 3s ease-out;
  }

  .message-fading {
    opacity: 0.3;
    transform: scale(0.95);
  }

  /* Notification Toast Styles */
  .notification-toast {
    @apply p-4 rounded-lg shadow-lg backdrop-blur-sm;
    @apply border border-gray-600/50;
    background: rgba(31, 41, 55, 0.95);
  }

  .notification-warning {
    @apply bg-gradient-to-r from-amber-900/50 to-orange-900/50;
    @apply border-amber-700/50;
    background: rgba(146, 64, 14, 0.95);
  }

  .notification-error {
    @apply bg-gradient-to-r from-red-900/50 to-pink-900/50;
    @apply border-red-700/50;
    background: rgba(153, 27, 27, 0.95);
  }

  .notification-success {
    @apply bg-gradient-to-r from-green-900/50 to-emerald-900/50;
    @apply border-green-700/50;
    background: rgba(20, 83, 45, 0.95);
  }

  .notification-info {
    @apply bg-gradient-to-r from-blue-900/50 to-cyan-900/50;
    @apply border-blue-700/50;
    background: rgba(30, 58, 138, 0.95);
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .text-gradient {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
}
</style>
