/** @type {import('next').NextConfig} */
const nextConfig = {
  //TODO: remover quando tirar do github pages //Só usar export estático em produção (GitHub Pages)
  ...(process.env.NODE_ENV === 'production' && {
    output: 'export',
    trailingSlash: true,
    assetPrefix: '/leaderpro',
    basePath: '/leaderpro',
  }),
  images: {
    unoptimized: true
  },
}

module.exports = nextConfig