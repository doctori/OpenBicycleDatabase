<template>
<div>
    <v-navigation-drawer app>
      <v-container>
        <v-row class="text-center">
          <v-col>
            <v-img
              :src="logo"
              />
            <router-link
            :to="{ name: 'Home'}"
            >
              <h4 class="display-3 font-weight-bold">
                OCD
              </h4>
            </router-link>
          </v-col>
          <v-col >
            <p class="subheading font-weight-regular">
              {{ $t("messages.welcome_nav") }}

            </p>
          </v-col>

          <v-col
            class="mb-5"
            cols="12"
          >
          <router-link :to="{name:'Standards'}">
            <h2 class="headline font-weight-bold mb-5">
              Standards
            </h2>
          </router-link>

          </v-col>

          <v-col
            class="mb-5"
            cols="12"
          >
          <router-link
            :to="{name: 'Components'}">
            <h2 class="headline font-weight-bold mb-5">
              {{ $t('messages.components')}}
            </h2>
          </router-link>
          </v-col>

          <v-col
            class="mb-5"
            cols="12"
          >
          <router-link
            :to="{name:'Brands'}"
            >
            <h2 class="headline font-weight-bold mb-5">
              Brands
            </h2>
          </router-link>
          </v-col>
        </v-row>
      </v-container>
    </v-navigation-drawer>
</div>
</template>

<script>
import logo from '../assets/logo.png';
export default {
  name: 'Navigation',

  data() {    
    return{
      logo,
    }}
}
</script>
