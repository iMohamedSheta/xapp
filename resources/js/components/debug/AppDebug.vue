<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watchEffect, computed } from 'vue'
import { GripVertical, Clock, Database, Server, User, Zap, Component, BugPlayIcon } from 'lucide-vue-next'
import { usePage } from '@inertiajs/vue3'
import { AppPageProps } from '../../types/index'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { ChevronDown, ChevronRight } from 'lucide-vue-next'
import Button from '../ui/button/Button.vue'
import { JsonViewer } from '@anilkumarthakur/vue3-json-viewer';
import '@anilkumarthakur/vue3-json-viewer/styles.css';


const tabs = [
  { id: 'timeline', label: 'Timeline' },
  { id: 'request', label: 'Request' },
  { id: 'queries', label: 'Queries' + ' (' + (usePage<AppPageProps>().props.debug?.queries?.length || 0) + ')' },
  { id: 'components', label: 'Components' },
  // { id: 'logs', label: 'Logs' },
  { id: 'raw', label: 'Raw' },
]

const activeTab = ref<string | null>(null)
const isDebugbarVisible = ref(true)

const page = usePage<AppPageProps>()
const debug = computed(() => page.props.debug)

// Get all component instances from the current page
const componentsData = computed(() => {
  // Otherwise, try to detect from Vue tree
  const components: any[] = []
  const visited = new WeakSet()

  function traverseComponent(app: any, depth = 0) {
    if (!app || visited.has(app)) return

    visited.add(app)

    // Get component info
    const type = app.type
    const componentName = type?.name || type?.__name || type?.displayName || 'Anonymous'
    const file = type?.__file || ''

    // Only add if it has a name and it's not internal Vue components
    if (componentName && componentName !== 'Anonymous') {
      // Check if it's a page - must be in /Pages/ or /pages/ directory (not just named *Page)
      const isPage = file.includes('/Pages/') ||
        file.includes('/pages/') ||
        file.includes('\\Pages\\') ||
        file.includes('\\pages\\')

      components.push({
        name: componentName,
        type: isPage ? 'page' : 'component',
        depth,
        file: file || 'Unknown',
        props: app.props || {},
        propsCount: Object.keys(app.props || {}).length,
        expanded: false
      })
    }

    // Traverse children components
    if (app.subTree) {
      traverseVNode(app.subTree, depth)
    }

    // Traverse child component instances
    if (app.children) {
      Object.values(app.children).forEach((child: any) => {
        if (child && typeof child === 'object') {
          traverseComponent(child, depth + 1)
        }
      })
    }
  }

  function traverseVNode(vnode: any, depth: number) {
    if (!vnode) return

    // If this vnode has a component instance, traverse it
    if (vnode.component) {
      traverseComponent(vnode.component, depth + 1)
    }

    // Traverse children vnodes
    if (vnode.children) {
      if (Array.isArray(vnode.children)) {
        vnode.children.forEach((child: any) => {
          if (child && typeof child === 'object') {
            traverseVNode(child, depth)
          }
        })
      }
    }

    // Traverse dynamicChildren if present
    if (vnode.dynamicChildren) {
      vnode.dynamicChildren.forEach((child: any) => {
        traverseVNode(child, depth)
      })
    }
  }

  // Get the root component instance from Inertia/Vue
  try {
    // Try to get from Inertia's page component
    const inertiaApp = document.getElementById('app')
    if (inertiaApp && (inertiaApp as any).__vueParentComponent) {
      traverseComponent((inertiaApp as any).__vueParentComponent, 0)
    } else if (inertiaApp && (inertiaApp as any).__vue_app__) {
      const app = (inertiaApp as any).__vue_app__
      if (app._instance) {
        traverseComponent(app._instance, 0)
      }
    }
  } catch (e) {
    console.error('Error traversing components:', e)
  }

  // Remove duplicates based on name+file combination
  const uniqueComponents = components.filter((comp, index, self) =>
    index === self.findIndex((c) => c.name === comp.name && c.file === comp.file)
  )

  // Sort: pages first, then components
  return uniqueComponents.sort((a, b) => {
    if (a.type === 'page' && b.type === 'component') return -1
    if (a.type === 'component' && b.type === 'page') return 1
    return a.name.localeCompare(b.name)
  })
})

function toggleTab(id: string) {
  activeTab.value = activeTab.value === id ? null : id
}

function closeDebugbar() {
  activeTab.value = null
  localStorage.setItem('debugbar_visible', 'false')
  isDebugbarVisible.value = false
}

function openDebugbar() {
  isDebugbarVisible.value = true
  localStorage.setItem('debugbar_visible', 'true')
}

function handleClickOutside(event: MouseEvent) {
  const debugbar = document.getElementById('debugbar')
  if (!debugbar) return

  const target = event.target as Node
  const isInsideDebugbar = debugbar.contains(target)

  if (!isInsideDebugbar) {
    activeTab.value = null
  }
}

function convertToMs(duration: string): number {
  if (duration.endsWith('µs')) {
    return parseFloat(duration.replace('µs', '')) / 1000
  } else if (duration.endsWith('ms')) {
    return parseFloat(duration.replace('ms', ''))
  } else if (duration.endsWith('s')) {
    return parseFloat(duration.replace('s', '')) * 1000
  }
  return parseFloat(duration)
}

// Generate timeline data based on debug info
const timelineData = computed(() => {
  if (!debug.value?.request) return []

  const totalDuration = convertToMs(debug.value.request.action_duration)
  const queryCount = debug.value.queries?.length || 0

  // Calculate query total time
  const queryTotalTime = debug.value.queries?.reduce((total, query) => {
    return total + convertToMs(query.execDuration)
  }, 0) || 0

  const timeline = []

  if (queryCount > 0) {
    timeline.push({
      id: 'database',
      label: 'Database Queries',
      icon: Database,
      timestamp: '~2ms',
      duration: `${queryTotalTime.toFixed(2)}ms`,
      type: 'process',
      description: `${queryCount} quer${queryCount === 1 ? 'y' : 'ies'} executed`,
      color: 'bg-green-500'
    })
  }

  timeline.push(
    {
      id: 'controller',
      label: 'Controller Processing',
      icon: User,
      timestamp: `~${(2 + queryTotalTime).toFixed(1)}ms`,
      duration: `${Math.max(1, totalDuration - queryTotalTime - 5).toFixed(2)}ms`,
      type: 'process',
      description: 'Business logic execution',
      color: 'bg-orange-500'
    },
    {
      id: 'response',
      label: 'Response Generated',
      icon: Clock,
      timestamp: `${totalDuration.toFixed(2)}ms`,
      duration: '0ms',
      type: 'milestone',
      description: `Status: ${debug.value.request.status}`,
      color: debug.value.request.status >= 400 ? 'bg-red-500' :
        debug.value.request.status >= 300 ? 'bg-yellow-500' : 'bg-green-500'
    }
  )

  return timeline
})

// Resizing logic
const panelHeight = ref(240)
const minHeight = 100
const maxHeight = 1080

let startY = 0
let startHeight = 0
let isResizing = false

function onMouseDown(event: MouseEvent) {
  isResizing = true
  startY = event.clientY
  startHeight = panelHeight.value
  event.preventDefault()
}

function onMouseMove(event: MouseEvent) {
  if (!isResizing) return
  const dy = startY - event.clientY
  let newHeight = startHeight + dy
  if (newHeight < minHeight) newHeight = minHeight
  if (newHeight > maxHeight) newHeight = maxHeight
  panelHeight.value = newHeight
}

function onMouseUp() {
  isResizing = false
}

onMounted(() => {
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
  isDebugbarVisible.value = localStorage.getItem('debugbar_visible') === 'true' ? true : false
})

function highlightSql(sql: string) {
  return sql.replace(
    /\b(SELECT|FROM|WHERE|LIMIT|INSERT|INTO|VALUES|UPDATE|DELETE|JOIN|INNER|LEFT|RIGHT|FULL|OUTER|ON|AS|DISTINCT|GROUP|BY|ORDER|HAVING|UNION|ALL|EXCEPT|INTERSECT|CREATE|TABLE|PRIMARY|KEY|FOREIGN|NOT|NULL|DEFAULT|CHECK|CONSTRAINT|ALTER|ADD|DROP|TRUNCATE|INDEX|VIEW|SEQUENCE|TRIGGER|BEGIN|COMMIT|ROLLBACK|GRANT|REVOKE|CASCADE|AND|OR|IN|IS|BETWEEN|LIKE|ILIKE|EXISTS|CASE|WHEN|THEN|ELSE|END)\b/gi,
    '<span class="kw text-yellow-600">$1</span>'
  );
}

function getFileName(filePath: string) {
  if (!filePath || filePath === 'Unknown') return filePath
  return filePath.split('/').pop() || filePath
}

function getShortPath(filePath: string) {
  if (!filePath || filePath === 'Unknown') return filePath

  // Find the index of /js/ and return everything from there
  const jsIndex = filePath.indexOf('/js/')
  if (jsIndex !== -1) {
    return filePath.substring(jsIndex)
  }

  // Also check for \js\ (Windows path)
  const jsIndexWin = filePath.indexOf('\\js\\')
  if (jsIndexWin !== -1) {
    return filePath.substring(jsIndexWin).replace(/\\/g, '/')
  }

  return filePath
}
</script>

<template>
  <!-- Floating Debug Button (shows when debugbar is hidden) -->
  <transition name="fade">
    <button v-if="!isDebugbarVisible" @click="openDebugbar"
      class="fixed bottom-0 left-0 z-[9999] p-1 bg-[#fafafa] dark:bg-[#1a1a1a] hover:bg-[#e6e4e4] dark:hover:bg-[#252525] border-2 border-[#353535] shadow-lg  transition-all duration-200 hover:scale-110"
      title="Open Debugbar" aria-label="Open Debugbar">
      <BugPlayIcon class="size-6 text-red-500 dark:text-red-400" />
    </button>
  </transition>

  <div id="debugbar" dir="ltr" class="fixed bottom-0 left-0 right-0 z-[10] text-xs debugbar-font"
    v-show="isDebugbarVisible">
    <!-- Expanded panel -->
    <transition name="slide-fade" mode="out-in" appear>
      <div v-if="activeTab" key="panel"
        class="bg-[#d8d8d8] dark:bg-[#121212] dark:text-white border-t border-[#353535] border-x ">
        <!-- Resize handle -->
        <div @mousedown.prevent="onMouseDown" title="Drag to resize"
          class="h-[1px] w-full cursor-n-resize flex justify-center items-center bg-[#d8d8d8] dark:bg-gray-800 mb-2">
          <div class="bg-border z-10 flex h-4 w-3 items-center justify-center rounded-xs border"
            @mousedown.prevent="onMouseDown" title="Drag to resize">
            <GripVertical class="size-2.5" />
          </div>
        </div>

        <div class="flex justify-between mb-2 border-b border-[#353535]">
          <h2 class="font-semibold text-sm capitalize px-3 pb-2 ">{{ activeTab }}</h2>
          <button class="text-gray-400 hover:text-white px-4" @click="activeTab = null" aria-label="Close debug panel">
            ✕
          </button>
        </div>

        <div class="py-2 px-4 overflow-auto scrollbar-tailwind" :style="{ height: panelHeight + 'px' }">

          <template v-if="activeTab === 'timeline'">
            <div class="space-y-4">

              <!-- Timeline Header -->
              <div class="flex items-center justify-between mb-4 pb-2 border-b border-gray-300 dark:border-gray-600">
                <div class="flex items-center space-x-2">
                  <Clock class="size-4" />
                  <span class="font-semibold">Request Timeline</span>
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  Total: {{ debug?.request?.action_duration }}
                </div>
              </div>

              <!-- Timeline Scale -->
              <div class="mb-4">
                <div class="flex justify-between text-xs text-gray-500 dark:text-gray-400 mb-1">
                  <span>0ms</span>
                  <span>{{ Math.round(convertToMs(debug?.request?.action_duration || '0ms') / 2) }}ms</span>
                  <span>{{ debug?.request?.action_duration }}</span>
                </div>
                <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded relative">
                  <div
                    class="absolute inset-0 bg-gradient-to-r from-green-400 via-yellow-400 to-red-400 rounded opacity-30">
                  </div>
                </div>
              </div>

              <!-- Timeline Items -->
              <div class="space-y-3">
                <div v-for="(item, index) in timelineData" :key="item.id" class="timeline-item">
                  <!-- Item Header -->
                  <div class="flex items-center justify-between mb-1">
                    <div class="flex items-center space-x-2">
                      <component :is="item.icon" :class="['size-4', item.color.replace('bg-', 'text-')]" />
                      <span class="text-sm font-medium">{{ item.label }}</span>
                    </div>
                    <div class="flex items-center space-x-3 text-xs text-gray-500 dark:text-gray-400">
                      <span>{{ convertToMs(item.timestamp.replace('~', '')) }}ms</span>
                      <span v-if="item.duration !== '0ms'"
                        class="px-2 py-0.5 bg-gray-100 dark:bg-gray-700 rounded font-mono">
                        {{ item.duration }}
                      </span>
                    </div>
                  </div>

                  <!-- Timeline Bar -->
                  <div class="relative h-6 bg-gray-100 dark:bg-gray-800 rounded border overflow-hidden">
                    <!-- Background grid lines -->
                    <div class="absolute inset-0 flex">
                      <div class="flex-1 border-r border-gray-200 dark:border-gray-700"></div>
                      <div class="flex-1 border-r border-gray-200 dark:border-gray-700"></div>
                      <div class="flex-1"></div>
                    </div>

                    <!-- Progress Bar -->
                    <div v-if="item.duration !== '0ms'" :class="[
                      'absolute top-0 bottom-0 rounded flex items-center justify-center text-white text-xs font-medium',
                      item.color
                    ]" :style="{
                      right: (convertToMs(item.timestamp.replace('~', '')) / convertToMs(debug?.request?.action_duration || '0ms')) * 100 + '%',
                      width: Math.max(2, (convertToMs(item.duration) / convertToMs(debug?.request?.action_duration || '0ms')) * 100) + '%'
                    }">
                      <span v-if="convertToMs(item.duration) > 5" class="truncate px-1">
                        {{ item.duration }}
                      </span>
                    </div>

                    <!-- Milestone Marker -->
                    <div v-else class="absolute top-0 bottom-0 w-1 bg-current"
                      :class="item.color.replace('bg-', 'text-')" :style="{
                        right: (convertToMs(item.timestamp.replace('~', '')) / convertToMs(debug?.request?.action_duration || '0ms')) * 100 + '%'
                      }">
                    </div>
                  </div>

                  <!-- Item Description -->
                  <p class="text-xs text-gray-600 dark:text-gray-400 mt-1 ml-6">{{ item.description }}</p>
                </div>
              </div>

              <!-- Performance Summary -->
              <div class="mt-6 pt-4 border-t border-gray-300 dark:border-gray-600">
                <h4 class="font-semibold text-sm mb-3">Performance Summary</h4>
                <div class="grid grid-cols-2 gap-4 text-xs">
                  <div class="space-y-2">
                    <div class="flex justify-between items-center">
                      <span>Request Duration:</span>
                      <span class="font-mono px-2 py-1 rounded" :class="{
                        'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200': convertToMs(debug?.request?.action_duration || '0ms') <= 100,
                        'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200': convertToMs(debug?.request?.action_duration || '0ms') > 100 && convertToMs(debug?.request?.action_duration || '0ms') < 500,
                        'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200': convertToMs(debug?.request?.action_duration || '0ms') >= 500
                      }">
                        {{ debug?.request?.action_duration }}
                      </span>
                    </div>
                    <div class="flex justify-between items-center">
                      <span>Database Queries:</span>
                      <span
                        class="font-mono px-2 py-1 bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 rounded">
                        {{ debug?.queries?.length || 0 }}
                      </span>
                    </div>
                  </div>
                  <div class="space-y-2">
                    <div class="flex justify-between items-center">
                      <span>Memory Usage:</span>
                      <span class="font-mono px-2 py-1 rounded" :class="{
                        'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200': debug?.request?.memory_usage_for_request?.endsWith('KB'),
                        'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200': debug?.request?.memory_usage_for_request?.endsWith('MB') && !debug?.request?.memory_usage_for_request?.startsWith('1'),
                        'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200': debug?.request?.memory_usage_for_request?.includes('MB') && (debug?.request?.memory_usage_for_request?.startsWith('1') || parseInt(debug?.request?.memory_usage_for_request) > 10)
                      }">
                        {{ debug?.request?.memory_usage_for_request }}
                      </span>
                    </div>
                    <div class="flex justify-between items-center">
                      <span>Status Code:</span>
                      <span class="font-mono px-2 py-1 rounded" :class="{
                        'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200': (debug?.request?.status ?? 0) >= 200 && (debug?.request?.status ?? 0) < 300,
                        'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200': (debug?.request?.status ?? 0) >= 300 && (debug?.request?.status ?? 0) < 400,
                        'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200': (debug?.request?.status ?? 0) >= 400
                      }">
                        {{ debug?.request?.status }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="activeTab === 'components'">
            <div class="space-y-4">
              <!-- Components Header -->
              <div class="flex items-center justify-between mb-4 pb-2 border-b border-gray-300 dark:border-gray-600">
                <div class="flex items-center space-x-2">
                  <Component class="size-4" />
                  <span class="font-semibold">Rendered Components</span>
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  Total: {{ componentsData.length }}
                </div>
              </div>

              <!-- Component Stats -->
              <div class="grid grid-cols-2 gap-4 mb-4">
                <div class="bg-blue-100 dark:bg-blue-900/30 p-3 rounded">
                  <div class="text-xs text-gray-600 dark:text-gray-400">Pages</div>
                  <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">
                    {{componentsData.filter(c => c.type === 'page').length}}
                  </div>
                </div>
                <div class="bg-purple-100 dark:bg-purple-900/30 p-3 rounded">
                  <div class="text-xs text-gray-600 dark:text-gray-400">Components</div>
                  <div class="text-2xl font-bold text-purple-600 dark:text-purple-400">
                    {{componentsData.filter(c => c.type === 'component').length}}
                  </div>
                </div>
              </div>

              <!-- Components List -->
              <div class="space-y-2">
                <div v-for="(comp, idx) in componentsData" :key="idx"
                  class="bg-gray-100 dark:bg-gray-800 rounded p-3 border-l-4" :class="{
                    'border-blue-500': comp.type === 'page',
                    'border-purple-500': comp.type === 'component'
                  }">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-3">
                      <span class="px-2 py-0.5 text-xs font-semibold rounded" :class="{
                        'bg-blue-500 text-white': comp.type === 'page',
                        'bg-purple-500 text-white': comp.type === 'component'
                      }">
                        {{ comp.type === 'page' ? 'PAGE' : 'COMP' }}
                      </span>
                      <div class="flex flex-col">
                        <span class="font-semibold">{{ comp.name }}</span>
                        <code class="text-xs text-gray-600 dark:text-gray-400 mt-0.5">
                          {{ getShortPath(comp.file) }}
                        </code>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Empty State -->
                <div v-if="componentsData.length === 0" class="text-center py-8 text-gray-500 dark:text-gray-400">
                  <Component class="size-12 mx-auto mb-2 opacity-50" />
                  <p>No components detected</p>
                  <p class="text-xs mt-1">Components will appear here when the page renders</p>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="activeTab === 'queries'">
            <ul class="debug-queries">
              <li v-for="(q, i) in debug?.queries" :key="i" class="debug-query my-2">
                <Collapsible v-model:open="q.expanded">
                  <CollapsibleTrigger class="w-full">
                    <div
                      class="flex items-start justify-between gap-4 cursor-pointer hover:bg-black/20 p-2 rounded -m-2">
                      <div class="flex-1 text-left">
                        <code v-html="highlightSql(q.sql)"></code>
                        <div class="flex justify-between">
                          <span class="time">
                            build: {{ q.buildDuration }} | exec: {{ q.execDuration }}
                          </span>
                          <span class="text-gray-400">
                            {{ q.sourceFile }} | Line: {{ q.sourceLine }}
                          </span>
                        </div>
                      </div>

                      <div class="flex-shrink-0 p-1">
                        <ChevronDown v-if="q.expanded" class="h-4 w-4" />
                        <ChevronRight v-else class="h-4 w-4" />
                      </div>
                    </div>
                  </CollapsibleTrigger>

                  <CollapsibleContent class="mt-3 space-y-3 border-t border-gray-700 pt-3">
                    <!-- Raw SQL -->
                    <div v-if="q.rawSql">
                      <h4 class="text-xs font-semibold text-gray-400 mb-2">Raw SQL:</h4>
                      <pre
                        class="bg-black/30 p-3 rounded text-xs overflow-x-auto"><code v-html="highlightSql(q.rawSql)"></code></pre>
                    </div>

                    <!-- Bindings -->
                    <div v-if="q.bindings && q.bindings.length > 0">
                      <h4 class="text-xs font-semibold text-gray-400 mb-2">Bindings:</h4>
                      <div class="bg-black/30 p-3 rounded">
                        <ul class="space-y-1">
                          <li v-for="(binding, idx) in q.bindings" :key="idx" class="text-xs">
                            <span class="text-gray-500">[{{ idx }}]:</span>
                            <span class="ml-2 font-mono text-green-400">{{ binding }}</span>
                          </li>
                        </ul>
                      </div>
                    </div>

                    <!-- No bindings message -->
                    <div v-else-if="!q.rawSql">
                      <p class="text-xs text-gray-500 italic">No additional data available</p>
                    </div>
                  </CollapsibleContent>
                </Collapsible>
              </li>
            </ul>
          </template>

          <template v-else-if="activeTab === 'request'">
            <div v-if="debug?.request" class="space-y-4">
              <!-- Request Section -->
              <div class="bg-[#fafafa] dark:bg-[#1e1e1e] border border-[#353535] p-4 shadow-lg">
                <h3 class="font-semibold text-sm mb-3 text-gray-700 dark:text-gray-300">Request Information</h3>

                <!-- Method + Path -->
                <div class="flex items-center space-x-2 mb-3">
                  <span class="px-2 py-1 text-xs font-bold text-black bg-yellow-600">
                    {{ debug.request.method }}
                  </span>
                  <span class="text-sm font-medium break-all">
                    {{ debug.request.path }}
                    <span v-if="debug.request.query" class="text-gray-400">?{{ debug.request.query }}</span>
                  </span>
                  <span class="ml-auto px-2 py-0.5 text-xs font-semibold" :class="{
                    'bg-green-600': debug.request.status >= 200 && debug.request.status < 300,
                    'bg-yellow-600': debug.request.status >= 300 && debug.request.status < 400,
                    'bg-red-600': debug.request.status >= 400
                  }">
                    {{ debug.request.status }}
                  </span>
                </div>

                <!-- Request Details Grid -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs">
                  <div><span class="font-semibold ">Request ID: </span> {{ debug.request.id }}</div>
                  <div><span class="font-semibold ">Client IP: </span> {{ debug.request.client_ip }}</div>
                  <div><span class="font-semibold ">Host: </span> {{ debug.request.host }}</div>
                  <div><span class="font-semibold ">Protocol: </span> {{ debug.request.protocol }}</div>
                  <div><span class="font-semibold ">Referer: </span> {{ debug.request.referer || '-' }}</div>
                  <div>
                    <span class="font-semibold">Duration: </span>
                    <span :class="{
                      'text-green-500': convertToMs(debug.request.action_duration) <= 10,
                      'text-yellow-500': convertToMs(debug.request.action_duration) > 10 && convertToMs(debug.request.action_duration) < 60,
                      'text-red-500': convertToMs(debug.request.action_duration) >= 60
                    }">
                      {{ debug.request.action_duration }}
                    </span>
                  </div>
                  <div>
                    <span class="font-semibold ">Memory (Request): </span>
                    <span
                      :class="{ 'text-green-500': debug.request.memory_usage_for_request.endsWith('KB'), 'text-red-500': debug.request.memory_usage_for_request.endsWith('MB') }">
                      {{ debug.request.memory_usage_for_request }}
                    </span>
                  </div>
                  <div>
                    <span class="font-semibold ">Memory (App): </span>
                    <span
                      :class="{ 'text-green-500': debug.request.memory_usage.endsWith('MB'), 'text-yellow-500': debug.request.memory_usage.endsWith('GB') }">
                      {{ debug.request.memory_usage }}
                    </span>
                  </div>
                  <div class="col-span-2"><span class="font-semibold ">User-Agent: </span>
                    <span class="text-[11px]">
                      {{ debug.request.user_agent }}
                    </span>
                  </div>
                </div>
              </div>

              <!-- Route Section -->
              <div v-if="debug?.route" class="bg-[#fafafa] dark:bg-[#1e1e1e] border border-[#353535] p-4 shadow-lg">
                <h3 class="font-semibold text-sm mb-3 text-gray-700 dark:text-gray-300">Route Information</h3>

                <div class="space-y-3 text-xs">
                  <!-- Route URI and Method -->
                  <div class="flex items-center space-x-2">
                    <span class="px-2 py-1 text-xs font-bold text-white bg-blue-600">
                      {{ debug.route.method }}
                    </span>
                    <code class="text-sm font-medium">{{ debug.route.uri }}</code>
                  </div>

                  <!-- Handler and Action -->
                  <div class="grid grid-cols-1 gap-2">
                    <div>
                      <span class="font-semibold">Handler: </span>
                      <code class="text-blue-600 dark:text-blue-400">{{ debug.route.handler }}</code>
                    </div>
                    <div>
                      <span class="font-semibold">Action: </span>
                      <code class="text-purple-600 dark:text-purple-400">{{ debug.route.action }}</code>
                    </div>
                    <div>
                      <span class="font-semibold">Short Name: </span>
                      <span class="text-gray-700 dark:text-gray-300">{{ debug.route.short }}</span>
                    </div>
                  </div>

                  <!-- Middleware -->
                  <div v-if="debug.route.middleware && debug.route.middleware.length > 0">
                    <span class="font-semibold block mb-2">Middleware Stack ({{ debug.route.middleware.length
                    }}):</span>
                    <div class="bg-gray-100 dark:bg-gray-800 rounded p-3 space-y-1">
                      <div v-for="(mw, idx) in debug.route.middleware" :key="idx" class="flex items-start space-x-2">
                        <span class="text-gray-500 dark:text-gray-400 font-mono">{{ idx + 1 }}.</span>
                        <code class="text-[11px] text-gray-700 dark:text-gray-300 break-all">{{ mw }}</code>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="activeTab === 'logs'">
            <p>App started at 10:00</p>
            <p>User logged in: user@example.com</p>
          </template>

          <template v-else-if="activeTab === 'raw'">
            <details>
              <summary>Raw Debug JSON</summary>
              <JsonViewer :data="debug" :level="6" :expanded="false" :darkMode="true" />
            </details>

            <details>
              <summary>Raw Props JSON</summary>
              <JsonViewer :data="page.props" :level="6" :expanded="false" :darkMode="true" />
            </details>
          </template>
        </div>
      </div>
    </transition>

    <!-- Bottom bar -->
    <nav
      class="relative bg-[#fafafa] dark:bg-[#121212] border-t border-[#353535] flex items-center justify-center space-x-6 px-4 py-2 select-none">
      <Button @click="closeDebugbar"
        class="absolute left-0 bg-[#fafafa] dark:bg-[#1a1a1a] hover:bg-[#e6e4e4] dark:hover:bg-[#252525] shadow-lg  transition-all duration-200 hover:scale-110">
        <BugPlayIcon class="size-6 text-red-500 dark:text-red-400" />
      </Button>

      <button v-for="tab in tabs" :key="tab.id" @click="toggleTab(tab.id)" :class="[
        'py-1 px-2 rounded-sm hover:bg-[#e6e4e4] focus:outline-none dark:hover:bg-[#202020] dark:focus:outline-none',
        activeTab === tab.id
          ? 'dark:bg-[#222121] dark:text-white bg-[#d1d1d1] font-semibold'
          : 'dark:text-gray-400',
      ]" :aria-pressed="activeTab === tab.id">
        {{ tab.label }}
      </button>
      <button class="text-red-400 hover:text-red-500 absolute right-4 py-1 px-3 " @click="closeDebugbar"
        aria-label="Close debug panel">
        ✕
      </button>
    </nav>
  </div>
</template>

<style scoped>
.kw {
  font-weight: bold;
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.slide-fade-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

.debugbar-font,
.debugbar-font :deep(*) {
  font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace !important;
  --font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace !important;
  --font-sans: 'Cascadia Code', 'Fira Code', 'Consolas', monospace !important;
  --font-mono: 'Cascadia Code', 'Fira Code', 'Consolas', monospace !important;
  font-size: 12px !important;
}

.debug-queries,
.debug-queries :deep(*) {
  font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace !important;
  font-size: 12px !important;
}

.debug-query {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 0.5rem 0.75rem;
  border-bottom: 1px solid #333;
}

.debug-query code {
  display: block;
  white-space: pre-wrap;
}

.debug-query .time {
  font-size: 0.8rem;
  color: #999;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: scale(0.8);
}

.fade-leave-to {
  opacity: 0;
  transform: scale(0.8);
}
</style>