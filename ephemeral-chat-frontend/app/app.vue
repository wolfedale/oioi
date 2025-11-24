<template>
  <div class="min-h-screen bg-gray-900 flex flex-col">
    <NuxtRouteAnnouncer />

    <!-- Main Content Area -->
    <div class="flex-1 flex items-center justify-center p-4">
      <!-- Main Chat Container -->
      <div class="w-full max-w-md">
        <!-- Account Creation Section -->
        <div v-if="!currentUser" class="glass-dark rounded-lg p-8 mb-6 text-center animate-fade-in">
          <div class="mb-6">
            <h1 class="text-3xl font-bold text-gradient mb-2">oi.oi</h1>
            <p class="text-gray-400 text-sm">Messages that exist only in the moment</p>
          </div>

          <!-- Step 1: Create Account or Login -->
          <div v-if="!hasAccount" class="space-y-4">
            <!-- Toggle between Create/Login -->
            <div class="flex bg-gray-800 rounded-lg p-1 mb-4">
              <button
                @click="isLoginMode = false"
                :class="[
                  'flex-1 py-2 px-4 rounded-md transition-all duration-200 text-sm font-medium',
                  !isLoginMode ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-gray-200'
                ]"
              >
                Create Account
              </button>
              <button
                @click="isLoginMode = true"
                :class="[
                  'flex-1 py-2 px-4 rounded-md transition-all duration-200 text-sm font-medium',
                  isLoginMode ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-gray-200'
                ]"
              >
                Login
              </button>
            </div>

            <!-- Create Account Mode -->
            <div v-if="!isLoginMode">
              <div class="mb-4">
                <h3 class="text-lg font-medium text-white mb-2">Create Account</h3>
                <p class="text-gray-400 text-xs">Write a sentence to generate your unique ID</p>
              </div>
              <textarea
                v-model="sentence"
                placeholder="Write a unique sentence (minimum 10 characters)..."
                class="input-modern-dark font-medium text-center resize-none h-20"
                @keyup.ctrl.enter="createAccount"
              ></textarea>
              <button
                @click="createAccount"
                :disabled="!sentence.trim() || sentence.trim().length < 10"
                class="btn-primary-dark w-full disabled:opacity-50 disabled:cursor-not-allowed mt-4"
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
                class="input-modern-dark font-medium text-center mb-3 font-mono"
                maxlength="8"
              />
              <textarea
                v-model="sentence"
                placeholder="Your original sentence..."
                class="input-modern-dark font-medium text-center resize-none h-20"
                @keyup.ctrl.enter="loginAccount"
              ></textarea>
              <button
                @click="loginAccount"
                :disabled="!loginUserID.trim() || !sentence.trim() || sentence.trim().length < 10"
                class="btn-primary-dark w-full disabled:opacity-50 disabled:cursor-not-allowed mt-4"
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
            <div class="bg-gray-800/80 border border-gray-600 rounded-lg p-4 mb-4">
              <p class="text-2xl font-mono font-bold text-blue-400 tracking-widest">{{ generatedUserID }}</p>
            </div>
            <button
              @click="useAccount"
              class="btn-primary-dark w-full"
            >
              Start Chatting
            </button>
            <p class="text-gray-500 text-xs">Save this ID - you'll need your sentence to regenerate it</p>
          </div>
        </div>

        <!-- Chat Interface -->
        <div v-else class="glass-dark rounded-lg p-6 animate-slide-up">
          <!-- Header -->
          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center space-x-3">
              <div class="w-3 h-3 rounded-full bg-green-500 animate-pulse-slow"></div>
              <div>
                <h2 class="font-semibold text-white">{{ currentUser }}</h2>
                <p class="text-xs text-gray-400">Connected</p>
              </div>
            </div>

            <button
              @click="leaveChat"
              class="text-gray-400 hover:text-gray-200 transition-colors p-2 rounded-lg hover:bg-gray-700"
            >
              <LogOut class="w-4 h-4" />
            </button>
          </div>

          <!-- Message Container -->
          <div
            ref="messagesContainer"
            class="h-80 overflow-y-auto scroll-smooth mb-4 space-y-3"
          >
            <!-- System Messages -->
            <div v-if="messages.length === 0" class="text-center py-8">
              <MessageCircle class="w-12 h-12 text-gray-600 mx-auto mb-3" />
              <p class="text-gray-400 text-sm">No messages yet...</p>
              <p class="text-gray-500 text-xs mt-1">Start a conversation</p>
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
                class="flex justify-end"
              >
                <div class="chat-bubble-sent-dark">
                  <p class="text-sm">{{ message.content }}</p>
                  <p class="text-xs text-blue-200 mt-1 opacity-75">
                    {{ formatTime(message.timestamp) }}
                  </p>
                </div>
              </div>

              <div v-else class="flex justify-start">
                <div class="chat-bubble-received-dark">
                  <div class="flex items-center space-x-2 mb-1">
                    <div class="w-2 h-2 rounded-full bg-gradient-to-r from-blue-500 to-purple-600"></div>
                    <p class="text-xs font-medium text-gray-300">{{ message.from }}</p>
                  </div>
                  <p class="text-sm">{{ message.content }}</p>
                  <p class="text-xs text-gray-400 mt-1">
                    {{ formatTime(message.timestamp) }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Message Input -->
          <div class="flex space-x-3">
            <div class="flex-1 relative">
              <input
                v-model="newMessage"
                type="text"
                placeholder="Type an ephemeral message..."
                class="input-modern-dark pr-12"
                @keyup.enter="sendMessage"
                :disabled="!targetUser"
              />
              <Send
                class="absolute right-4 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400"
              />
            </div>
          </div>

          <!-- Target User Input -->
          <div class="mt-4 pt-4 border-t border-gray-700">
            <label class="block text-xs font-medium text-gray-300 mb-2">Send to:</label>
            <input
              v-model="targetUser"
              type="text"
              placeholder="Enter recipient username..."
              class="input-modern-dark text-sm"
            />
          </div>
        </div>

        <!-- Footer -->
        <div class="text-center mt-6 text-xs text-gray-500">
          <p>Messages disappear after reading • No history stored</p>
        </div>
      </div>
    </div>

    <!-- Copyright Footer - Fixed at bottom -->
    <div class="text-center pb-4 text-xs text-gray-600">
      <p>© 2024 REDGRIP.io. All rights reserved.</p>
      <p class="mt-1">Powered by REDGRIP Tier</p>
    </div>

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
import { MessageCircle, Send, LogOut, Plus, X } from 'lucide-vue-next'

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
let messagePollingInterval: NodeJS.Timeout | null = null
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

  // Start message polling
  startMessagePolling()

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
  if (messagePollingInterval) clearInterval(messagePollingInterval)
  if (uiUpdateInterval) clearInterval(uiUpdateInterval)

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

  // Initial ping
  await pingPresence()

  // Set up interval
  presenceInterval = setInterval(pingPresence, 15000) // Every 15 seconds
}

const startMessagePolling = () => {
  const pollMessages = async () => {
    try {
      const response = await fetch(`${config.public.apiBase}/message/receive?user_id=${currentUser.value}`)
      const data = await response.json()

      if (data.message) {
        // Add readAt timestamp
        const messageWithReadTime = {
          ...data.message,
          readAt: new Date().toISOString()
        }

        // Determine which conversation this message belongs to
        const conversationId = messageWithReadTime.from === currentUser.value
          ? messageWithReadTime.to
          : messageWithReadTime.from

        // Get or create conversation
        const conversation = getOrCreateConversation(conversationId)

        // Check for duplicate messages (by ID) in this conversation
        const existingMessage = conversation.messages.find(msg => msg.id === messageWithReadTime.id)
        if (existingMessage) {
          console.log(`Duplicate message detected, skipping: ${messageWithReadTime.id}`)
          return
        }

        // Add message to conversation
        conversation.messages.push(messageWithReadTime)
        conversation.lastMessage = messageWithReadTime
        conversation.lastActivity = new Date().toISOString()

        // Schedule removal after 10 seconds
        scheduleMessageRemoval(messageWithReadTime.id, conversationId)

        console.log(`Message received in conversation ${conversationId} and will disappear in 10 seconds: ${messageWithReadTime.id}`)

        // Auto-scroll to bottom if this is the active conversation
        if (activeConversationId.value === conversationId) {
          nextTick(() => {
            if (messagesContainer.value) {
              messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
            }
          })
        }
      }
    } catch (error) {
      console.error('Message polling failed:', error)
    }
  }

  // Poll every 1 second for better responsiveness
  messagePollingInterval = setInterval(pollMessages, 1000)
}


const sendMessage = async () => {
  if (!newMessage.value.trim() || !activeConversationId.value) return

  const targetUser = activeConversationId.value

  // Ensure sender is marked as online before sending
  try {
    await fetch(`${config.public.apiBase}/presence/ping`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: currentUser.value })
    })
  } catch (error) {
    console.error('Failed to ping presence before sending:', error)
  }

  try {
    const response = await fetch(`${config.public.apiBase}/message/send`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        from: currentUser.value,
        to: targetUser.value,
        content: newMessage.value.trim()
      })
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${data.error || 'Unknown error'}`)
    }

    // Add sent message immediately to conversation for instant feedback
    const sentMessage = {
      id: data.message_id,
      from: currentUser.value,
      to: targetUser,
      content: newMessage.value.trim(),
      timestamp: data.timestamp,
      expires_at: new Date(Date.now() + 60000).toISOString(),
      readAt: new Date().toISOString()
    }

    // Get the conversation and add the message
    const conversation = getOrCreateConversation(targetUser)
    conversation.messages.push(sentMessage)
    conversation.lastMessage = sentMessage
    conversation.lastActivity = new Date().toISOString()

    scheduleMessageRemoval(sentMessage.id, targetUser)
    console.log(`Message sent to ${targetUser}: ${sentMessage.id}`)

    // Clear input
    newMessage.value = ''

    // Auto-scroll to bottom
    nextTick(() => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    })

  } catch (error: any) {
    console.error('Send message error:', error)

    // Generic error message to maintain privacy
    showNotification('Message could not be delivered. The recipient may not be available.', 'warning')
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
  if (messagePollingInterval) clearInterval(messagePollingInterval)
  if (uiUpdateInterval) clearInterval(uiUpdateInterval)

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
    @apply bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50;
    @apply text-slate-900 antialiased;
    min-height: 100vh;
  }

  @media (prefers-color-scheme: dark) {
    body {
      @apply bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900;
      @apply text-slate-100;
    }
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
