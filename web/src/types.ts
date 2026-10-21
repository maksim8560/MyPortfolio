// Зеркало Go-структур из server/content.go.
// Держим в синхроне вручную: контракт маленький и меняется редко.

export interface Link {
  label: string
  href: string
  external?: boolean
}

export interface Fact {
  key: string
  value: string
}

export interface Stat {
  value: string
  label: string
}

export interface Project {
  id: string
  name: string
  kind: string
  tagline: string
  summary: string
  bullets: string[]
  tags: string[]
  demo?: string
  demoText?: string
  year: string
  status: string
  featured: boolean
}

export interface ApproachItem {
  title: string
  body: string
}

export interface ContactItem {
  label: string
  value: string
  href?: string
  copy?: boolean
  copied?: string
}

export interface Section {
  kicker: string
  title: string
  lead: string
}

export interface Content {
  lang: string
  meta: { title: string; description: string }
  nav: Link[]
  hero: {
    kicker: string
    title1: string
    title2: string
    lead: string
    primary: Link
    secondary: Link
    facts: Fact[]
    note: string
  }
  stats: Stat[]
  work: Section
  projects: Project[]
  approach: Section & { items: ApproachItem[] }
  contact: Section & { items: ContactItem[]; note: string }
  footer: { note: string; built: string }
}

export interface Health {
  status: string
  version: string
  uptime: number
  runtime: string
}
