import type { CreatePostPayload, PostItem, PostListQuery } from './post'

export interface PostFormModel extends CreatePostPayload {
  id: number
}

export type PostPageQuery = PostListQuery
export type PostPageItem = PostItem
