package fproto_doc

import (
	"sort"

	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto/fdep"
)

// Sort type for filter
type SortType int

const (
	ST_NONE                SortType = iota // No sort
	ST_FILEPATH_ALIAS_NAME                 // proto file path + alias + name
	ST_ALIAS_NAME                          // alias + name
	ST_NAME                                // name
)

// Dependencies to filter
type FilterDepType int

const (
	DT_ALL      FilterDepType = iota // All dependency types
	DT_OWN                           // Only own dependencies
	DT_IMPORTED                      // Only imported dependencies
)

// Filter struct
type GetFilter struct {
	SortType      SortType
	FilterDepType FilterDepType
	FilePaths     []string
}

func NewGetFilter(sortType SortType, filterDepType FilterDepType) *GetFilter {
	return &GetFilter{
		SortType:      sortType,
		FilterDepType: filterDepType,
	}
}

func (gf *GetFilter) SetFilePaths(filePaths []string) *GetFilter {
	gf.FilePaths = filePaths
	return gf
}

// Doc generator struct
type Helper struct {
	dep *fdep.Dep
}

// Creates a new doc generator
func NewHelper(dep *fdep.Dep) *Helper {
	return &Helper{
		dep: dep,
	}
}

// Get a list of all enums using the filter
func (g *Helper) GetEnumList(filter *GetFilter) []*fdep.DepType {
	return g.genList(filter, func(pfile *fproto.ProtoFile) []fproto.FProtoElement {
		return pfile.CollectEnums()
	})
}

// Get a list of all messages using the filter
func (g *Helper) GetMessageList(filter *GetFilter) []*fdep.DepType {
	return g.genList(filter, func(pfile *fproto.ProtoFile) []fproto.FProtoElement {
		return pfile.CollectMessages()
	})
}

// Get a list of all services using the filter
func (g *Helper) GetServiceList(filter *GetFilter) []*fdep.DepType {
	return g.genList(filter, func(pfile *fproto.ProtoFile) []fproto.FProtoElement {
		return pfile.CollectServices()
	})
}

// Get a list of all services using the filter
func (g *Helper) GetOneOfFieldList(fields []fproto.FieldElementTag) []fproto.FieldElementTag {
	var ret []fproto.FieldElementTag
	for _, fld := range fields {
		if oofld, is_oneof := fld.(*fproto.OneofFieldElement); is_oneof {
			ret = append(ret, oofld)
			ret = append(ret, g.GetOneOfFieldList(oofld.Fields)...)
		}
	}
	return ret
}

// Get a list of fields sorted by name
func (g *Helper) SortedFieldList(fields []fproto.FieldElementTag) []fproto.FieldElementTag {
	var ret []fproto.FieldElementTag

	collect := make(map[string]fproto.FieldElementTag)
	for _, fld := range fields {
		collect[fld.FieldName()] = fld
	}

	var skeys []string
	for k, _ := range collect {
		skeys = append(skeys, k)
	}

	sort.Sort(sort.StringSlice(skeys))

	for _, k := range skeys {
		ret = append(ret, collect[k])
	}

	return ret
}

// Get a list of fields sorted by tag
func (g *Helper) SortedByTagFieldList(fields []fproto.FieldElementTag) []fproto.FieldElementTag {
	var ret []fproto.FieldElementTag

	collect := make(map[int]fproto.FieldElementTag)
	for _, fld := range fields {
		collect[fld.FirstFieldTag()] = fld
	}

	var skeys []int
	for k, _ := range collect {
		skeys = append(skeys, k)
	}

	sort.Sort(sort.IntSlice(skeys))

	for _, k := range skeys {
		ret = append(ret, collect[k])
	}

	return ret
}

// Get a list of RPCs sorted by name
func (g *Helper) SortedRPCList(rpcs []*fproto.RPCElement) []*fproto.RPCElement {
	var ret []*fproto.RPCElement

	collect := make(map[string]*fproto.RPCElement)
	for _, fld := range rpcs {
		collect[fld.Name] = fld
	}

	var skeys []string
	for k, _ := range collect {
		skeys = append(skeys, k)
	}

	sort.Sort(sort.StringSlice(skeys))

	for _, k := range skeys {
		ret = append(ret, collect[k])
	}

	return ret
}

// Get the list of files sorted by name
func (g *Helper) SortedFileList(filterDepType FilterDepType) []string {
	var skeys []string
	for k, f := range g.dep.Files {
		include := true
		switch filterDepType {
		case DT_OWN:
			include = f.DepType == fdep.DepType_Own
		case DT_IMPORTED:
			include = f.DepType == fdep.DepType_Imported
		}

		if include {
			skeys = append(skeys, k)
		}
	}

	sort.Sort(sort.StringSlice(skeys))

	return skeys
}

// Get the list of packages sorted by name
func (g *Helper) SortedPackageList(filterDepType FilterDepType) []string {
	var skeys []string
	for k, fl := range g.dep.Packages {
		include := true
		if filterDepType != DT_ALL {
			// check if any file of the package matches the filter
			include = false
			for _, f := range fl {
				switch filterDepType {
				case DT_OWN:
					include = g.dep.Files[f].DepType == fdep.DepType_Own
				case DT_IMPORTED:
					include = g.dep.Files[f].DepType == fdep.DepType_Imported
				}
			}
			if include {
				break
			}
		}

		if include {
			skeys = append(skeys, k)
		}
	}

	sort.Sort(sort.StringSlice(skeys))

	return skeys
}

// Internal list generator
func (g *Helper) genList(filter *GetFilter, pffunc func(pfile *fproto.ProtoFile) []fproto.FProtoElement) []*fdep.DepType {
	collect := make(map[string]*fdep.DepType)
	var ret []*fdep.DepType

	for _, f := range g.dep.Files {
		include := true
		switch filter.FilterDepType {
		case DT_OWN:
			include = f.DepType == fdep.DepType_Own
		case DT_IMPORTED:
			include = f.DepType == fdep.DepType_Imported
		}

		if include && len(filter.FilePaths) > 0 {
			include = false
			for _, fp := range filter.FilePaths {
				if fp == f.FilePath {
					include = true
					break
				}
			}
		}

		if include {
			for _, e := range pffunc(f.ProtoFile) {
				dt := fdep.NewDepTypeFromElement(f, e)
				if filter.SortType == ST_NONE {
					ret = append(ret, dt)
				} else {
					collect[g.sortTypeValue(filter.SortType, dt)] = dt
				}
			}
		}
	}

	if filter.SortType == ST_NONE {
		return ret
	}

	return g.sortDepType(collect)
}

func (g *Helper) sortDepType(m map[string]*fdep.DepType) []*fdep.DepType {
	var skeys []string
	for k, _ := range m {
		skeys = append(skeys, k)
	}

	sort.Sort(sort.StringSlice(skeys))

	var ret []*fdep.DepType
	for _, k := range skeys {
		ret = append(ret, m[k])
	}
	return ret
}

func (g *Helper) sortTypeValue(sortType SortType, dp *fdep.DepType) string {
	switch sortType {
	case ST_ALIAS_NAME:
		return dp.Alias + "." + dp.Name
	case ST_NAME:
		return dp.Name
	default:
		return dp.FileDep.FilePath + "." + dp.Alias + "." + dp.Name
	}
}
