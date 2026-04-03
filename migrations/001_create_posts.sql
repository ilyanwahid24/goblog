CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    excerpt TEXT DEFAULT '',
    content TEXT NOT NULL,
    published BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
CREATE INDEX IF NOT EXISTS idx_posts_published ON posts(published);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);

-- Seed some sample posts
INSERT INTO posts (title, slug, excerpt, content, published) VALUES
(
    'Welcome to My Blog',
    'welcome-to-my-blog',
    'This is the first post on this blog. Welcome aboard!',
    '<p>Hello and welcome to my blog! This is a simple blog built with <strong>Go</strong> and <strong>PostgreSQL</strong>.</p>
<p>I built this as a clean, minimal blogging platform that focuses on what matters most — the content. No unnecessary bloat, no heavy JavaScript frameworks, just clean server-rendered pages that load fast.</p>
<h2>Features</h2>
<ul>
<li>Server-side rendered with Go templates</li>
<li>PostgreSQL for reliable data storage</li>
<li>Clean, responsive design</li>
<li>Simple admin panel for managing posts</li>
</ul>
<p>Stay tuned for more posts!</p>',
    true
),
(
    'Getting Started with Go',
    'getting-started-with-go',
    'A beginner-friendly guide to the Go programming language.',
    '<p>Go is a statically typed, compiled programming language designed at Google. It''s known for its simplicity, efficiency, and excellent support for concurrent programming.</p>
<h2>Why Go?</h2>
<p>Go was designed to address the challenges of building large-scale software systems. Here are some reasons why developers love Go:</p>
<ul>
<li><strong>Simplicity:</strong> Go has a clean, minimal syntax that''s easy to learn</li>
<li><strong>Performance:</strong> Compiled to native code, Go programs run fast</li>
<li><strong>Concurrency:</strong> Goroutines make concurrent programming straightforward</li>
<li><strong>Standard Library:</strong> Rich standard library for web servers, crypto, and more</li>
</ul>
<h2>Hello World</h2>
<pre><code>package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}</code></pre>
<p>That''s all you need to get started. Go is a fantastic language for building web applications, CLI tools, and microservices.</p>',
    true
),
(
    'Why PostgreSQL is Great',
    'why-postgresql-is-great',
    'Exploring the strengths of PostgreSQL as a database choice.',
    '<p>PostgreSQL is one of the most advanced open-source relational database systems. It has earned a reputation for reliability, data integrity, and correctness.</p>
<h2>Key Strengths</h2>
<ul>
<li><strong>ACID Compliance:</strong> Full transactional support ensures data integrity</li>
<li><strong>Extensibility:</strong> Custom types, operators, and functions</li>
<li><strong>JSON Support:</strong> First-class JSON and JSONB data types</li>
<li><strong>Full Text Search:</strong> Built-in full-text search capabilities</li>
<li><strong>Reliability:</strong> Battle-tested in production at massive scale</li>
</ul>
<p>For this blog, PostgreSQL provides a solid foundation for storing and querying posts efficiently. Its indexing capabilities ensure that even with thousands of posts, queries remain fast.</p>',
    true
)
ON CONFLICT (slug) DO NOTHING;
