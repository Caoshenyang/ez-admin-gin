import type { PostFormModel, PostPageQuery } from '../types/post-page'
import { PostStatus, type PostItem } from '../types/post'

export function defaultPostQuery(): PostPageQuery {
  return {
    keyword: '',
    status: 0,
  }
}

export function defaultPostFormModel(): PostFormModel {
  return {
    id: 0,
    code: '',
    name: '',
    sort: 0,
    status: PostStatus.Enabled,
    remark: '',
  }
}

export function toPostFormModel(post: PostItem): PostFormModel {
  return {
    id: post.id,
    code: post.code,
    name: post.name,
    sort: post.sort,
    status: post.status,
    remark: post.remark,
  }
}

export function buildPostPayload(formModel: PostFormModel) {
  return {
    code: formModel.code.trim(),
    name: formModel.name.trim(),
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}
