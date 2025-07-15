'use client'

import { Calendar, Building2, User, MapPin } from 'lucide-react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Person } from '@/lib/types'
import { formatTimeAgoWithoutSuffix, formatLastOneOnOne } from '@/lib/utils/dates'
import { getInitials } from '@/lib/utils/names'

interface PersonCardProps {
  person: Person
  className?: string
  onClick?: () => void
}

export function PersonCard({ 
  person, 
  className = '',
  onClick
}: PersonCardProps) {
  return (
    <Card 
      className={`hover:shadow-md transition-all hover:scale-[1.02] cursor-pointer ${className}`}
      onClick={onClick}
    >
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-3">
            <Avatar className="h-12 w-12">
              <AvatarImage src={person.avatar} alt={person.name} />
              <AvatarFallback className="bg-primary/10 text-primary">
                {getInitials(person.name)}
              </AvatarFallback>
            </Avatar>
            
            <div className="flex-1 min-w-0">
              <h3 className="font-semibold truncate">
                {person.name}
              </h3>
              <p className="text-sm text-muted-foreground truncate">
                {person.position || person.role || 'Cargo não informado'}
              </p>
              
              <div className="flex flex-col gap-1 mt-1">
                {person.department && (
                  <div className="flex items-center gap-1">
                    <Building2 className="h-3 w-3 text-muted-foreground" />
                    <span className="text-xs text-muted-foreground">
                      {person.department}
                    </span>
                  </div>
                )}
                {person.primaryAddress && (person.primaryAddress.city || person.primaryAddress.state) && (
                  <div className="flex items-center gap-1">
                    <MapPin className="h-3 w-3 text-muted-foreground" />
                    <span className="text-xs text-muted-foreground">
                      {person.primaryAddress.city && person.primaryAddress.state
                        ? `${person.primaryAddress.city}, ${person.primaryAddress.state}`
                        : person.primaryAddress.city || person.primaryAddress.state}
                    </span>
                  </div>
                )}
              </div>
            </div>
          </div>

          {person.hasKids && (
            <Badge variant="outline" className="text-xs">
              Tem filhos
            </Badge>
          )}
        </div>
      </CardHeader>

      <CardContent className="space-y-3">
        {/* Last 1:1 - More relevant info */}
        <div className="flex items-center gap-2">
          <Calendar className="h-4 w-4 text-muted-foreground" />
          <div className="flex-1">
            <p className="text-sm text-muted-foreground">
              Último 1:1:
            </p>
            <p className="text-sm font-medium">
              {formatLastOneOnOne(person.lastOneOnOneDate)}
            </p>
          </div>
        </div>

        {/* Time at company */}
        <div className="flex items-center gap-2">
          <User className="h-4 w-4 text-muted-foreground" />
          <div>
            <p className="text-sm text-muted-foreground">
              Na empresa há{' '}
              <span className="font-medium">
                {person.startDate ? formatTimeAgoWithoutSuffix(person.startDate) : 'não informado'}
              </span>
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}