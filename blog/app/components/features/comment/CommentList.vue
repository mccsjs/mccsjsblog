<script setup lang="ts">
import type { Comment } from '@@/types/comment';
import CommentItem from './CommentItem.vue';

interface FlatComment {
  comment: Comment;
  depth: number;
}

interface Props {
  comments: FlatComment[];
}

const props = defineProps<Props>();

const groupedComments = computed(() => {
  const groups: Array<{
    parent: FlatComment;
    replies: FlatComment[];
  }> = [];

  let currentGroup: { parent: FlatComment; replies: FlatComment[] } | null = null;

  props.comments.forEach(item => {
    if (item.depth === 0) {
      currentGroup = {
        parent: item,
        replies: [],
      };
      groups.push(currentGroup);
    } else {
      if (currentGroup) {
        currentGroup.replies.push(item);
      }
    }
  });

  return groups;
});
</script>

<template>
  <div class="comments-list">
    <div v-for="group in groupedComments" :key="group.parent.comment.id" class="comment-card">
      <CommentItem :comment="group.parent.comment" :depth="group.parent.depth">
        <template v-if="group.replies.length > 0" #replies>
          <div class="replies-list">
            <CommentItem
              v-for="reply in group.replies"
              :key="reply.comment.id"
              :comment="reply.comment"
              :depth="reply.depth"
            />
          </div>
        </template>
      </CommentItem>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.comments-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-card {
  background: var(--mccsjs-card-bg);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 3px 8px 6px rgba(7, 17, 27, 0.05);
  transition: 0.3s;

  [data-theme='dark'] & {
    box-shadow: none;
    border: 1px solid var(--mccsjs-border);
  }

  &:hover {
    box-shadow: 0 3px 8px 6px rgba(7, 17, 27, 0.09);

    [data-theme='dark'] & {
      box-shadow: none;
    }
  }
}

@media screen and (max-width: 768px) {
  .comment-card {
    padding: 12px;
  }
}
</style>
