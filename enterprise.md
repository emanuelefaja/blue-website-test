# Enterprise-Only Transition Plan

## Overview
Transitioning Blue from B2B SaaS focused on SMEs to an Enterprise-only strategy. This change reflects that 80-90% of growth comes from enterprise customers, primarily through the partner program.

## Key Changes Required

### 1. Navbar Updates
- [✅] Replace current buttons (Login, Contact, Sign Up) with: Support, Login, Book Demo
- [✅] Update `components/topbar.html` for desktop and mobile navigation
- [✅] Update translations in `translations/common.json` for all 16 languages
- [✅] Fix responsive text wrapping issue with whitespace-nowrap

### 2. Create Book Demo Page
- [ ] Create new page at `pages/book-demo.html` 
- [ ] Use `contact/sales.html` as template
- [ ] Include demo scheduling widget (Calendly/HubSpot)
- [ ] Add enterprise customer logos
- [ ] Focus on enterprise value propositions

### 3. Pricing Page Replacement
- [✅] Remove all pricing content from `pages/pricing.html` (page deleted)
- [✅] Add redirect from `/pricing` to `/demo` in `data/redirects.json`
- [✅] Remove pricing from navigation in `data/nav.json`

### 4. Documentation Updates

#### Billing Documentation (`content/en/docs/11.billing.md`)
- [✅] Remove pricing from left sidebar navigation
- [✅] Remove specific pricing ($7/month, $70/year per user)
- [✅] Remove 7-day free trial references
- [✅] Replace with enterprise billing information:
  - Custom pricing
  - Invoice-based billing
  - Annual contracts
  - Volume discounts
- [✅] Update for all 15 translated versions (auto-translated with cmd/translate-docs)

#### Other Documentation
- [ ] Remove free trial mentions
- [ ] Remove self-service signup instructions
- [ ] Add enterprise onboarding process documentation

### 5. Homepage Updates (`pages/index.html`)
- [ ] Remove "$7/user. No surprises. Save 90% vs enterprise tools"
- [✅] Replace "Get Started Free" CTAs with "Book a Demo" (Updated dual-button-cta to use /demo)
- [ ] Update hero section for enterprise positioning
- [ ] Emphasize enterprise customers: CBRE, Cisco, Deloitte, FedEx, Manulife, Red Hat

### 6. CTA Component Updates
- [✅] Updated `components/dual-button-cta.html`:
  - Changed logged-in redirect from app.blue.cc to /demo
  - Updated button text from "Launch App" to "Book Demo"
  - Changed title for logged-in users to "Ready to see Blue for your enterprise?"
- [✅] Updated key pages to use /demo instead of app.blue.cc:
  - Homepage (`pages/index.html`)
  - Sales contact page (`pages/contact/sales.html`)
- [✅] Verified `components/cta-single-button.html` (no changes needed)

### 7. Content Updates Across Site

#### Files requiring updates (39 insights/blog posts with pricing/SME mentions):
- [ ] `content/en/insights/pricing-change.md`
- [ ] `content/en/insights/ideal-customer-profile.md`
- [ ] `content/en/insights/cash-is-king.md`
- [ ] `content/en/insights/2021-annual-letter.md`
- [ ] `content/en/insights/2022-annual-letter.md`
- [ ] Other insights files (see full list below)

#### Global changes needed:
- [ ] Replace all "Sign Up" buttons with "Book Demo" or "Contact Sales"
- [ ] Remove mentions of free trials
- [ ] Remove self-service onboarding references
- [ ] Update meta descriptions for enterprise focus

### 7. FAQ and Support Updates
- [ ] Update `content/en/insights/security-faq.md`
- [ ] Update `content/en/insights/ai-faq.md`
- [ ] Remove pricing-related FAQs
- [ ] Add enterprise-focused FAQs:
  - Security and compliance
  - Deployment options
  - Integration capabilities
  - Support SLAs

### 8. Translation Updates
- [ ] Update button text in all 16 languages:
  - English (en)
  - Chinese Simplified (zh)
  - Spanish (es)
  - French (fr)
  - German (de)
  - Japanese (ja)
  - Portuguese (pt)
  - Russian (ru)
  - Korean (ko)
  - Italian (it)
  - Indonesian (id)
  - Dutch (nl)
  - Polish (pl)
  - Chinese Traditional (zh-TW)
  - Swedish (sv)
  - Khmer (km)

### 9. Additional Enhancements
- [ ] Create enterprise features page highlighting:
  - SSO/SAML authentication
  - Advanced audit trails
  - Custom user roles and permissions
  - API access and webhooks
  - Dedicated infrastructure options
- [ ] Add enterprise case studies
- [ ] Update value propositions to emphasize:
  - Enterprise-grade security
  - Compliance certifications
  - Custom integrations
  - Dedicated support
  - Scalability

## Implementation Priority
1. **High Priority** (Do First):
   - Navbar changes
   - Create book demo page
   - Pricing page redirect
   - Homepage CTA updates

2. **Medium Priority** (Do Second):
   - Billing documentation update
   - Remove free trial mentions
   - Update translations

3. **Lower Priority** (Do Third):
   - Update all blog posts
   - Add enterprise case studies
   - Create enterprise features page

## Testing Checklist
- [✅] All redirects work correctly (pricing redirects to /demo)
- [✅] No broken links to pricing page (removed from nav, redirects added)
- [ ] Book demo form/widget functions properly
- [ ] Navigation works in all languages
- [✅] No remaining mentions of $7/$70 pricing in billing docs
- [✅] CTA components updated to use "Book Demo" instead of "Sign Up"
- [ ] Mobile navigation updated correctly

## Files to Review/Update

### Critical Files:
- `components/topbar.html` - Navigation bar
- `pages/pricing.html` - Pricing page (to redirect)
- `pages/index.html` - Homepage
- `content/en/docs/11.billing.md` - Billing documentation
- `data/redirects.json` - Redirect configuration
- `translations/common.json` - Common UI translations
- `translations/pricing.json` - Pricing page translations

### Supporting Files:
- All files in `content/*/insights/` directories
- All solution pages in `pages/solutions/`
- Platform feature pages in `pages/platform/features/`

## Success Metrics
- Clear enterprise positioning throughout site
- No confusion about pricing model
- Smooth demo booking process
- Consistent messaging across all languages
- Professional enterprise-grade appearance

## Notes
- Maintain high polish expected from enterprise software
- Ensure messaging aligns with companies like Palantir, Stripe, OpenAI, Linear, Figma
- Keep implementation simple and maintainable for single developer
- Partner program should be prominently featured as primary growth channel