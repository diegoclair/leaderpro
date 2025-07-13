export async function generateStaticParams() {
  // Static generation disabled since people are now dynamic from API
  // Dynamic routes are handled at runtime via [id] dynamic segment
  return []
}

export default function ProfileLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return <>{children}</>
}