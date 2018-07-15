import { Comparator, State } from 'clarity-angular';

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
