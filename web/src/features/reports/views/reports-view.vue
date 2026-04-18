<script setup lang="ts">
import { Card, Button, Input, Heading, Text, Icon } from "@/components/ui";
</script>

<template>
  <div class="max-w-7xl mx-auto px-6 lg:px-10 py-10">
    <div class="mb-8 animate-fade-up">
      <Heading as="h1" size="3xl">Reports</Heading>
      <Text as="p" size="sm" muted mt="1">Spending trends and financial insights</Text>
    </div>

    <!-- Date range filter -->
    <div class="flex flex-col sm:flex-row gap-4 mb-8 animate-fade-up delay-1 items-end">
      <Input label="Start Date" type="date" />
      <Input label="End Date" type="date" />
      <Button variant="secondary">Apply</Button>
    </div>

    <Card>
      <div class="p-12 text-center animate-fade-up delay-2">
        <Icon name="bar-chart-2" size="xl" class="mx-auto text-[var(--rule)] mb-4" />
        <Heading as="h3" size="lg" mb="1">No report data</Heading>
        <Text as="p" size="sm" muted>
          Add transactions to see spending breakdowns and trends.
        </Text>
      </div>
    </Card>
  </div>
</template>