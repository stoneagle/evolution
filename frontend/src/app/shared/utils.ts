import { Comparator, State } from 'clarity-angular';
import { Serializable }      from '../model/base/serializable';
import { Basic, PageSet, Sort } from '../model/base/basic';
import { FilterType }        from './const';

export class CustomComparator<T> implements Comparator<T> {
  fieldName: string;
  type: string;

  constructor(fieldName: string, type: string) {
    this.fieldName = fieldName;
    this.type = type;
  }

  compare(a: { [key: string]: any | any[] }, b: { [key: string]: any | any[] }) {
    let comp = 0;
    if (a && b) {
      let fieldA = a[this.fieldName];
      let fieldB = b[this.fieldName];
      switch (this.type) {
        case "number":
          comp = fieldB - fieldA;
          break;
        case "date":
          comp = new Date(fieldB).getTime() - new Date(fieldA).getTime();
          break;
      }
    }
    return comp;
  }
}

export function reloadState(currentState: State, currentPageSet: PageSet): State {
  let state: State = currentState;
  if (!state) {
      state = {
          page: {}
      };
  }
  state.page.from = 0;
  state.page.to = currentPageSet.Size - 1;
  state.page.size = currentPageSet.Size;
  return state;
}

export function deleteState(currentPageSet: PageSet, currentState: State, deleteNum: number) {
  let total: number = currentPageSet.Count - deleteNum;
  if (total <= 0) { 
    return null; 
  }

  let totalPages: number = Math.ceil(total / currentPageSet.Size);
  let targetPageNumber: number = currentPageSet.Current;

  if (currentPageSet.Current > totalPages) {
      targetPageNumber = totalPages;
  }

  let state: State = currentState;
  if (!state) {
      state = { page: {} };
  }
  state.page.size = currentPageSet.Size;
  state.page.from = (targetPageNumber - 1) * currentPageSet.Size;
  state.page.to = targetPageNumber * currentPageSet.Size - 1;
  return state;
}

export function loadPageFilterSort<T extends Basic >(resource: T, state: State): T {
  if (!resource) {
    return resource;
  }

  let pageNumber = 0;
  if (!state || !state.page) {
    pageNumber = 1;
  } else {
    pageNumber = Math.ceil((state.page.to + 1) / state.page.size);
    if (pageNumber <= 0) {
      pageNumber = 1;
    }
  }
  if (!resource.Page) {
    resource.Page = new PageSet();
  }
  resource.Page.Size = state.page.size;
  resource.Page.Current = pageNumber;

  if (state && (state.filters != undefined) && state.filters.length > 0) {
    state.filters.forEach((filter: any) => {
      if (filter.ftype != undefined) {
        switch(filter.ftype) {
          case FilterType.Checkbox:
            break;
          case FilterType.Radio:
            resource[filter.filterParamName] = filter.selectedItem.key;
            break;
        }
      } else {
        let key = filter.property;
        let value = filter.value;
        resource[key] = value;
      }
    });
  }

  if (state && (state.sort != undefined)) {
    if (!resource.Sort) {
      resource.Sort = new Sort();
    }
    if (typeof(state.sort.by) == "string") {
      resource.Sort.By = state.sort.by 
    } else {
      resource.Sort.By = state.sort.by["fieldName"]
    }
    resource.Sort.Reverse = state.sort.reverse;
  }
  return resource;
}

export function doFiltering<T extends { [key: string]: any | any[] }>(items: T[], state: any): T[] {
    if (!items || items.length === 0) {
        return items;
    }

    if (!state || !state.filters || state.filters.length === 0) {
        return items;
    }

    state.filters.forEach((filter: {
        property: string;
        value: string;
    }) => {
        items = items.filter(item => regexpFilter(filter["value"], item[filter["property"]]));
    });

    return items;
}

export function doSorting<T extends { [key: string]: any | any[] }>(items: T[], state: State): T[] {
    if (!items || items.length === 0) {
        return items;
    }
    if (!state || !state.sort) {
        return items;
    }

    return items.sort((a: T, b: T) => {
        let comp: number = 0;
        if (typeof state.sort.by !== "string") {
            comp = state.sort.by.compare(a, b);
        } else {
            let propA = a[state.sort.by.toString()], propB = b[state.sort.by.toString()];
            if (typeof propA === "string") {
                comp = propA.localeCompare(propB);
            } else {
                if (propA > propB) {
                    comp = 1;
                } else if (propA < propB) {
                    comp = -1;
                }
            }
        }

        if (state.sort.reverse) {
            comp = -comp;
        }

        return comp;
    });
}

export function regexpFilter(terms: string, testedValue: any): boolean {
    let reg = new RegExp('.*' + terms + '.*', 'i');
    return reg.test(testedValue);
}
