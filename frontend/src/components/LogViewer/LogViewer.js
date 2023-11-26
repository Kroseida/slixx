import {defineComponent} from 'vue'

export default defineComponent({
  name: 'LogViewer',
  data() {
    return {
      displayType: true,
      logsConsole: ''
    }
  },
  props: {
    logs: {
      type: Array,
    }
  },
  watch: {
    logs() {
      this.logsConsole = this.logs.map((log) => {
        return `[${log.loggedAt}] [${log.level.toUpperCase()}] ${log.message}`
      }).join('\n')
    }
  },
  mounted() {
    this.logsConsole = this.logs.map((log) => {
      return `[${log.loggedAt}] [${log.level.toUpperCase()}] ${log.message}`
    }).join('\n')
  },
  methods: {
    classFromLevel(level) {
      if (level === 'error') {
        return 'bg-red-5'
      } else if (level === 'warn') {
        return 'bg-amber-5'
      } else if (level === 'info') {
        return 'bg-green-4'
      } else if (level === 'debug') {
        return 'bg-grey-2'
      }
      return 'bg-green-3'
    },
    downloadLogs() {
      const element = document.createElement('a')
      const file = new Blob([this.logsConsole], {type: 'text/plain'})
      element.href = URL.createObjectURL(file)
      element.download = 'logs.txt'
      document.body.appendChild(element) // Required for this to work in FireFox
      element.click()
    }
  }
})
