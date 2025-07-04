import { mockPeople } from '@/lib/data/mockData'

export async function generateStaticParams() {
  return mockPeople.map((person) => ({
    id: person.id,
  }))
}