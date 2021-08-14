<template>
  <div class="about">
    <h1 v-on:click="getUserList">{{ info }}</h1>
  </div>
</template>

<script>
import axios from "axios";
export default {
  name: "About",
  data() {
    return {
      info: this.getUserList(),
    };
  },
  methods: {
    getUserList() {
      axios
          .get("/try/ajax/json_demo.json")
          .then((response) => (this.info = response))
          .catch(function (error) {
            console.log(error);
          });
    },
  },
  mounted() {
    this.getUserList()
  }
}
</script>
