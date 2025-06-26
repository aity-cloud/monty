<script>
import Main from '../components/Main';
import { choosePageLayout } from '../utils/standalone';
import Alarm from '../components/Alarm';

export default {
  components: { Alarm, Main },
  layout:     choosePageLayout(),
  head() {
    return { title: 'Monty - Alarm' };
  }
};
</script>
<template>
  <Main>
    <Alarm />
  </Main>
</template>
