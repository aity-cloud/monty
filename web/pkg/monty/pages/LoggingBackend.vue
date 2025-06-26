<script>
import Main from '../components/Main';
import { choosePageLayout } from '../utils/standalone';
import LoggingBackend from '../components/LoggingBackend';

export default {
  components: { LoggingBackend, Main },
  layout:     choosePageLayout(),
  head() {
    return { title: 'Monty - Logging' };
  }
};
</script>
<template>
  <Main>
    <LoggingBackend />
  </Main>
</template>
