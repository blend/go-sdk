/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import "html/template"

// ViewCacheOption is an option for ViewCache.
type ViewCacheOption func(*ViewCache) error

// OptViewCachePaths sets the view cache paths.
func OptViewCachePaths(paths ...string) ViewCacheOption {
	return func(vc *ViewCache) error { vc.Paths = append(vc.Paths, paths...); return nil }
}

// OptViewCacheLiterals sets the view cache literals.
func OptViewCacheLiterals(literals ...string) ViewCacheOption {
	return func(vc *ViewCache) error { vc.Literals = append(vc.Literals, literals...); return nil }
}

// OptViewCacheFuncMap sets the view cache func maps.
func OptViewCacheFuncMap(funcMap template.FuncMap) ViewCacheOption {
	return func(vc *ViewCache) error { vc.FuncMap = funcMap; return nil }
}

// OptViewCacheFunc adds a func to the view func map.
func OptViewCacheFunc(name string, viewFunc interface{}) ViewCacheOption {
	return func(vc *ViewCache) error {
		if vc.FuncMap == nil {
			vc.FuncMap = make(template.FuncMap)
		}
		vc.FuncMap[name] = viewFunc
		return nil
	}
}

// OptViewCacheLiveReload adds a func to the view func map.
func OptViewCacheLiveReload(liveReload bool) ViewCacheOption {
	return func(vc *ViewCache) error { vc.LiveReload = liveReload; return nil }
}

// OptViewCacheInternalErrorTemplateName sets the internal error template name.
func OptViewCacheInternalErrorTemplateName(name string) ViewCacheOption {
	return func(vc *ViewCache) error { vc.InternalErrorTemplateName = name; return nil }
}

// OptViewCacheBadRequestTemplateName sets the bad request template name.
func OptViewCacheBadRequestTemplateName(name string) ViewCacheOption {
	return func(vc *ViewCache) error { vc.BadRequestTemplateName = name; return nil }
}

// OptViewCacheNotFoundTemplateName sets the not found template name.
func OptViewCacheNotFoundTemplateName(name string) ViewCacheOption {
	return func(vc *ViewCache) error { vc.NotFoundTemplateName = name; return nil }
}

// OptViewCacheNotAuthorizedTemplateName sets the not authorized template name.
func OptViewCacheNotAuthorizedTemplateName(name string) ViewCacheOption {
	return func(vc *ViewCache) error { vc.NotAuthorizedTemplateName = name; return nil }
}

// OptViewCacheStatusTemplateName sets the status template name.
func OptViewCacheStatusTemplateName(name string) ViewCacheOption {
	return func(vc *ViewCache) error { vc.StatusTemplateName = name; return nil }
}

// OptViewCacheConfig sets options based on a config.
func OptViewCacheConfig(cfg *ViewCacheConfig) ViewCacheOption {
	return func(vc *ViewCache) error {
		vc.Paths = cfg.Paths
		vc.LiveReload = cfg.LiveReload
		vc.InternalErrorTemplateName = cfg.InternalErrorTemplateNameOrDefault()
		vc.BadRequestTemplateName = cfg.BadRequestTemplateNameOrDefault()
		vc.NotFoundTemplateName = cfg.NotFoundTemplateNameOrDefault()
		vc.NotAuthorizedTemplateName = cfg.NotAuthorizedTemplateNameOrDefault()
		vc.StatusTemplateName = cfg.StatusTemplateNameOrDefault()
		return nil
	}
}
