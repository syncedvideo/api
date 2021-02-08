<template>
  <div>
    <button
      @click="show = !show"
      class="flex justify-center items-center leading-none"
    >
      <span class="">Settings</span>
    </button>

    <div
      v-if="show"
      class="fixed inset-0 flex justify-center items-center"
      style="background:rgba(0,0,0,.9)"
      @click.self="show = false"
    >
      <div>
        <input
          v-model="newColor"
          @keydown.enter="colorHandler"
          type="text"
          placeholder="Change color"
        />
        <input
          v-model="newUsername"
          @keydown.enter="usernameHandler"
          type="text"
          placeholder="Change username"
        />
      </div>
      <!-- <div>
        <div>Connected users</div>
        <ul>
          <li
            v-for="connection of connections"
            :key="connection.user.id"
            v-text="
              connection.user.username
                ? connection.user.username
                : connection.user.id
            "
            :style="{ color: connection.user.chatColor }"
          ></li>
        </ul>
      </div> -->
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'

export default defineComponent({
  name: 'AppRoomSettings',

  setup() {
    const show = ref(false)

    // const connections: ComputedRef<ConnectionMap> = computed(() => {
    //   return room.state.value.data.connectionHub.connections
    // })

    // const currentUser: ComputedRef<UserDto> = computed(() => {
    //   return room.state.value.user
    // })

    const newUsername = ref('')
    function usernameHandler() {
      // if (newUsername.value) {
      //   room.sendClientUsername(newUsername.value)
      //   newUsername.value = ''
      // }
    }

    const newColor = ref('')
    function colorHandler() {
      // if (newColor.value) {
      //   room.sendClientColor(newColor.value)
      //   newColor.value = ''
      // }
    }

    return {
      show,
      // connections,
      // currentUser,
      newUsername,
      usernameHandler,
      newColor,
      colorHandler
    }
  }
})
</script>
