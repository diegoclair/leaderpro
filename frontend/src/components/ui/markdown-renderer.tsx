'use client'

import React, { memo } from 'react'
import ReactMarkdown, { Components } from 'react-markdown'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark } from 'react-syntax-highlighter/dist/esm/styles/prism'
import remarkGfm from 'remark-gfm'
import { cn } from '@/lib/utils'

interface CodeProps {
  inline?: boolean
  className?: string
  children?: React.ReactNode
}

interface MarkdownRendererProps {
  content: string
  className?: string
}

// Memoized component for performance - prevents re-parsing on every token update
export const MarkdownRenderer = memo(function MarkdownRenderer({
  content,
  className
}: MarkdownRendererProps) {
  return (
    <div className={cn("prose prose-sm dark:prose-invert max-w-none", className)}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          // Code blocks with syntax highlighting
          code({ inline, className, children, ...props }: CodeProps & any) {
            const match = /language-(\w+)/.exec(className || '')
            const language = match ? match[1] : ''
            
            return !inline && match ? (
              <SyntaxHighlighter
                style={oneDark as any}
                language={language}
                PreTag="div"
                customStyle={{
                  margin: '1rem 0',
                  borderRadius: '8px',
                  fontSize: '0.875rem'
                }}
                {...props}
              >
                {String(children).replace(/\n$/, '')}
              </SyntaxHighlighter>
            ) : (
              <code 
                className={cn(
                  "relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm",
                  className
                )}
                {...props}
              >
                {children}
              </code>
            )
          },
          
          // Custom paragraph styling
          p({ children }: { children: React.ReactNode }) {
            return <p className="leading-7 [&:not(:first-child)]:mt-6">{children}</p>
          },
          
          // Custom heading styling
          h1({ children }: { children: React.ReactNode }) {
            return (
              <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
                {children}
              </h1>
            )
          },
          
          h2({ children }: { children: React.ReactNode }) {
            return (
              <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
                {children}
              </h2>
            )
          },
          
          h3({ children }: { children: React.ReactNode }) {
            return (
              <h3 className="scroll-m-20 text-2xl font-semibold tracking-tight">
                {children}
              </h3>
            )
          },
          
          // Custom list styling
          ul({ children }: { children: React.ReactNode }) {
            return <ul className="my-6 ml-6 list-disc [&>li]:mt-2">{children}</ul>
          },
          
          ol({ children }: { children: React.ReactNode }) {
            return <ol className="my-6 ml-6 list-decimal [&>li]:mt-2">{children}</ol>
          },
          
          // Custom blockquote styling
          blockquote({ children }: { children: React.ReactNode }) {
            return (
              <blockquote className="mt-6 border-l-2 pl-6 italic">
                {children}
              </blockquote>
            )
          },
          
          // Custom table styling
          table({ children }: { children: React.ReactNode }) {
            return (
              <div className="my-6 w-full overflow-y-auto">
                <table className="w-full">{children}</table>
              </div>
            )
          },
          
          tr({ children }: { children: React.ReactNode }) {
            return <tr className="m-0 border-t p-0 even:bg-muted">{children}</tr>
          },
          
          th({ children }: { children: React.ReactNode }) {
            return (
              <th className="border px-4 py-2 text-left font-bold [&[align=center]]:text-center [&[align=right]]:text-right">
                {children}
              </th>
            )
          },
          
          td({ children }: { children: React.ReactNode }) {
            return (
              <td className="border px-4 py-2 text-left [&[align=center]]:text-center [&[align=right]]:text-right">
                {children}
              </td>
            )
          },
        } as Components}
      >
        {content}
      </ReactMarkdown>
    </div>
  )
}, (prevProps, nextProps) => {
  // Custom comparison for memoization - only re-render if content actually changed
  return prevProps.content === nextProps.content && prevProps.className === nextProps.className
})