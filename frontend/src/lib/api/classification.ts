import { fetchAPI, authHeaders } from './core';

export interface ClassificationScheme {
  id: number;
  code: string;
  name_en: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  description?: string;
  icon?: string;
  status: 'ACTIVE' | 'SUSPENDED' | 'ARCHIVED' | 'DRAFT';
  version: number;
  revision: number;
  is_system: boolean;
  rules?: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface ClassificationConcept {
  id: number;
  code: string;
  canonical_name: string;
  name_en: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  description?: string;
  concept_type: 'PRODUCT' | 'SERVICE' | 'MATERIAL' | 'EQUIPMENT' | 'PROPERTY' | 'EMPLOYMENT' | 'ORGANIZATION' | 'OTHER';
  status: 'ACTIVE' | 'DRAFT' | 'REVIEW' | 'ARCHIVED' | 'RETIRED';
  brand_id?: number;
  metadata?: Record<string, any>;
  version: number;
  created_at: string;
  updated_at: string;
}

export interface ClassificationPlacement {
  id: number;
  scheme_id: number;
  concept_id: number;
  parent_placement_id?: number;
  role: 'ROOT' | 'SUPERFAMILY' | 'CATEGORY' | 'SUBCATEGORY' | 'LEAF' | 'OTHER';
  display_name_override?: string;
  slug: string;
  path: string;
  depth: number;
  display_order: number;
  is_leaf: boolean;
  is_active: boolean;
  status: 'DRAFT' | 'REVIEW' | 'APPROVED' | 'PUBLISHED' | 'ARCHIVED';
  created_at: string;
  updated_at: string;
  concept?: ClassificationConcept;
}

export interface PlacementTreeNode {
  id: number;
  parent_placement_id?: number;
  scheme_id: number;
  concept_id: number;
  role: string;
  slug: string;
  path: string;
  depth: number;
  display_order: number;
  is_leaf: boolean;
  is_active: boolean;
  name: string;
  name_en?: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  description?: string;
  icon?: string;
  concept_type: string;
  product_count: number;
  children: PlacementTreeNode[];
}

export interface ClassificationRelationshipType {
  code: string;
  name: string;
  description?: string;
  is_symmetric: boolean;
  inverse_code?: string;
}

export interface ClassificationRelationship {
  id: number;
  source_concept_id: number;
  target_concept_id: number;
  relationship_type_code: string;
  context_scheme_id?: number;
  status: string;
  metadata?: Record<string, any>;
  created_at: string;
  source_concept?: ClassificationConcept;
  target_concept?: ClassificationConcept;
}

export interface ClassificationAttributeRule {
  id: number;
  placement_id?: number;
  concept_id?: number;
  attribute_id: number;
  is_required: boolean;
  is_inherited: boolean;
  display_order: number;
  override_validation_rules?: Record<string, any>;
  override_unit?: string;
}

export interface ClassificationAttributeDTO {
  attribute_id: number;
  name_en: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  slug: string;
  data_type: string;
  validation_rules?: Record<string, any>;
  unit?: string;
  is_required: boolean;
  is_filterable: boolean;
  is_searchable: boolean;
  is_inherited: boolean;
  display_order: number;
}

export interface ImpactAnalysisResult {
  entity_id: number;
  entity_type: string;
  affected_placements_count: number;
  affected_listings_count: number;
  affected_descendants_count: number;
  affected_relationships: number;
  can_safely_execute: boolean;
  warnings: string[];
  details?: Record<string, any>;
}

export interface UniversalSearchResult {
  query: string;
  concepts: ClassificationConcept[];
  placements: PlacementTreeNode[];
  relationships?: ClassificationRelationship[];
  total_results: number;
}

export interface ClassificationAuditLog {
  id: number;
  entity_type: string;
  entity_id: number;
  action: string;
  performed_by?: number;
  before_state?: Record<string, any>;
  after_state?: Record<string, any>;
  reason?: string;
  created_at: string;
}

export const classificationApi = {
  // Public classification endpoints
  getSchemes: (status = 'ACTIVE') =>
    fetchAPI<ClassificationScheme[]>(`/classification/schemes?status=${encodeURIComponent(status)}`).then(r => r ?? []),

  getScheme: (code: string) =>
    fetchAPI<ClassificationScheme>(`/classification/schemes/${encodeURIComponent(code)}`),

  getTree: (schemeCode = 'SYS_PROD', lang = 'en', maxDepth = 4) =>
    fetchAPI<PlacementTreeNode[]>(`/classification/schemes/${encodeURIComponent(schemeCode)}/tree?lang=${encodeURIComponent(lang)}&depth=${maxDepth}`).then(r => r ?? []),

  getRoots: (schemeCode = 'SYS_PROD', lang = 'en') =>
    fetchAPI<PlacementTreeNode[]>(`/classification/schemes/${encodeURIComponent(schemeCode)}/roots?lang=${encodeURIComponent(lang)}`).then(r => r ?? []),

  getChildren: (placementId: number, lang = 'en') =>
    fetchAPI<PlacementTreeNode[]>(`/classification/placements/${placementId}/children?lang=${encodeURIComponent(lang)}`).then(r => r ?? []),

  getPlacement: (id: number, lang = 'en') =>
    fetchAPI<ClassificationPlacement>(`/classification/placements/${id}?lang=${encodeURIComponent(lang)}`),

  getBreadcrumbs: (placementId: number, lang = 'en') =>
    fetchAPI<any[]>(`/classification/placements/${placementId}/breadcrumbs?lang=${encodeURIComponent(lang)}`).then(r => r ?? []),

  getEffectiveAttributes: (placementId: number) =>
    fetchAPI<ClassificationAttributeDTO[]>(`/classification/placements/${placementId}/attributes`).then(r => r ?? []),

  getConcept: (id: number) =>
    fetchAPI<ClassificationConcept>(`/classification/concepts/${id}`),

  getRelationships: (conceptId: number, relType = '') => {
    const q = relType ? `?type=${encodeURIComponent(relType)}` : '';
    return fetchAPI<ClassificationRelationship[]>(`/classification/concepts/${conceptId}/relationships${q}`).then(r => r ?? []);
  },

  searchUniversal: (query: string, schemeCode = '', lang = 'en', limit = 20) => {
    const params = new URLSearchParams({
      q: query,
      lang,
      limit: String(limit),
    });
    if (schemeCode) params.set('scheme', schemeCode);
    return fetchAPI<UniversalSearchResult>(`/classification/search?${params.toString()}`);
  },

  checkDuplicateConcepts: (name: string, code = '') => {
    const params = new URLSearchParams({ name });
    if (code) params.set('code', code);
    return fetchAPI<ClassificationConcept[]>(`/classification/concepts/check-duplicate?${params.toString()}`).then(r => r ?? []);
  },

  // Admin classification endpoints
  adminCreateScheme: (data: Partial<ClassificationScheme>, token: string) =>
    fetchAPI<ClassificationScheme>('/admin/classification/schemes', {
      method: 'POST',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  adminUpdateScheme: (id: number, data: Partial<ClassificationScheme>, token: string) =>
    fetchAPI<ClassificationScheme>(`/admin/classification/schemes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  adminCreateConcept: (data: Partial<ClassificationConcept>, token: string) =>
    fetchAPI<ClassificationConcept>('/admin/classification/concepts', {
      method: 'POST',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  adminUpdateConcept: (id: number, data: Partial<ClassificationConcept>, token: string) =>
    fetchAPI<ClassificationConcept>(`/admin/classification/concepts/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  adminCreatePlacement: (data: Partial<ClassificationPlacement>, token: string) =>
    fetchAPI<ClassificationPlacement>('/admin/classification/placements', {
      method: 'POST',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  adminUpdatePlacement: (id: number, data: Partial<ClassificationPlacement>, token: string) =>
    fetchAPI<ClassificationPlacement>(`/admin/classification/placements/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  adminMovePlacement: (id: number, newParentPlacementId: number, token: string) =>
    fetchAPI<any>(`/admin/classification/placements/${id}/move`, {
      method: 'POST',
      body: JSON.stringify({ new_parent_placement_id: newParentPlacementId }),
      headers: authHeaders(token),
    }),

  adminDeletePlacement: (id: number, token: string) =>
    fetchAPI<any>(`/admin/classification/placements/${id}`, {
      method: 'DELETE',
      headers: authHeaders(token),
    }),

  adminAnalyzeImpact: (entityType: string, entityId: number, token: string) =>
    fetchAPI<ImpactAnalysisResult>(`/admin/classification/impact?type=${encodeURIComponent(entityType)}&id=${entityId}`, {
      headers: authHeaders(token),
    }),

  adminCreateRelationship: (data: Partial<ClassificationRelationship>, token: string) =>
    fetchAPI<any>('/admin/classification/relationships', {
      method: 'POST',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  adminDeleteRelationship: (id: number, token: string) =>
    fetchAPI<any>(`/admin/classification/relationships/${id}`, {
      method: 'DELETE',
      headers: authHeaders(token),
    }),

  adminSaveAttributeRule: (data: Partial<ClassificationAttributeRule>, token: string) =>
    fetchAPI<any>('/admin/classification/attributes/rules', {
      method: 'POST',
      body: JSON.stringify(data),
      headers: authHeaders(token),
    }),

  adminDeleteAttributeRule: (id: number, token: string) =>
    fetchAPI<any>(`/admin/classification/attributes/rules/${id}`, {
      method: 'DELETE',
      headers: authHeaders(token),
    }),

  adminListAuditLogs: (entityType?: string, entityId?: number, token?: string) => {
    const params = new URLSearchParams();
    if (entityType) params.set('type', entityType);
    if (entityId) params.set('id', String(entityId));
    return fetchAPI<ClassificationAuditLog[]>(`/admin/classification/audit-logs?${params.toString()}`, {
      headers: token ? authHeaders(token) : undefined,
    }).then(r => r ?? []);
  },
};

