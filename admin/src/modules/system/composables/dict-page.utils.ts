import type {
  DictItem,
  DictItemListQuery,
  DictStatus,
  DictTypeItem,
  DictTypeListQuery,
} from '../types/dict'
import type {
  DictItemCreatePayload,
  DictItemFormModel,
  DictItemUpdatePayload,
  DictTypeCreatePayload,
  DictTypeFormModel,
  DictTypeUpdatePayload,
} from '../types/dict-page'

import { DictStatus as DictStatusValue } from '../types/dict'

export function defaultDictTypeQuery(): DictTypeListQuery {
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    status: 0,
  }
}

export function defaultDictItemQuery(typeID = 0): DictItemListQuery {
  return {
    page: 1,
    page_size: 10,
    type_id: typeID,
    keyword: '',
    status: 0,
  }
}

function normalizeDictStatus(status: DictStatus | 0 | undefined) {
  return status === 0 ? undefined : status
}

export function normalizeDictTypeQuery(query: DictTypeListQuery): DictTypeListQuery {
  return {
    ...query,
    keyword: query.keyword?.trim() || undefined,
    status: normalizeDictStatus(query.status),
  }
}

export function normalizeDictItemQuery(
  query: DictItemListQuery,
  typeID: number,
): DictItemListQuery {
  return {
    ...query,
    type_id: typeID,
    keyword: query.keyword?.trim() || undefined,
    status: normalizeDictStatus(query.status),
  }
}

export function defaultDictTypeFormModel(): DictTypeFormModel {
  return {
    id: 0,
    code: '',
    name: '',
    sort: 10,
    status: DictStatusValue.Enabled,
    remark: '',
  }
}

export function defaultDictItemFormModel(typeID: number): DictItemFormModel {
  return {
    id: 0,
    type_id: typeID,
    item_key: '',
    label: '',
    value: '',
    tag_type: '',
    sort: 10,
    status: DictStatusValue.Enabled,
    remark: '',
  }
}

export function toDictTypeFormModel(item: DictTypeItem): DictTypeFormModel {
  return {
    id: item.id,
    code: item.code,
    name: item.name,
    sort: item.sort,
    status: item.status,
    remark: item.remark,
  }
}

export function toDictItemFormModel(item: DictItem): DictItemFormModel {
  return {
    id: item.id,
    type_id: item.type_id,
    item_key: item.item_key,
    label: item.label,
    value: item.value,
    tag_type: item.tag_type,
    sort: item.sort,
    status: item.status,
    remark: item.remark,
  }
}

export function buildDictTypeCreatePayload(formModel: DictTypeFormModel): DictTypeCreatePayload {
  return {
    code: formModel.code.trim(),
    name: formModel.name.trim(),
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}

export function buildDictTypeUpdatePayload(formModel: DictTypeFormModel): DictTypeUpdatePayload {
  return {
    name: formModel.name.trim(),
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}

export function buildDictItemCreatePayload(
  typeID: number,
  formModel: DictItemFormModel,
): DictItemCreatePayload {
  return {
    type_id: typeID,
    item_key: formModel.item_key.trim(),
    label: formModel.label.trim(),
    value: formModel.value.trim(),
    tag_type: formModel.tag_type.trim(),
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}

export function buildDictItemUpdatePayload(formModel: DictItemFormModel): DictItemUpdatePayload {
  return {
    label: formModel.label.trim(),
    value: formModel.value.trim(),
    tag_type: formModel.tag_type.trim(),
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}
