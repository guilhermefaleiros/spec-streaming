import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'
import { Menu } from 'lucide-react'
import { Link } from 'react-router-dom'
import { links } from './app-sidebar'

export function MobileNavDrawer() {
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button variant="outline" aria-label="Open navigation">
          <Menu aria-hidden="true" />
          Menu
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="w-72">
        <SheetHeader>
          <SheetTitle>Navigation</SheetTitle>
          <SheetDescription>Open a section of the app.</SheetDescription>
        </SheetHeader>
        <nav aria-label="Primary" className="flex flex-col gap-2">
          {links.map((link) => (
            <SheetClose key={link.to} asChild>
              <Link
                to={link.to}
                className="rounded-md px-3 py-2 text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
              >
                {link.label}
              </Link>
            </SheetClose>
          ))}
        </nav>
      </SheetContent>
    </Sheet>
  )
}
