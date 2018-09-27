import { Component, Input } from '@angular/core';
import { Filter }           from 'clarity-angular';
import { Subject }          from 'rxjs/Subject';
import { FilterType }       from '../const';

@Component({
  moduleId: 'share-filter-checkbox',
  selector: 'share-filter-checkbox',
  template: `
    <div>
      <ul style="list-style:none">
        <li *ngFor="let item of items">
          <clr-checkbox
            (change)="onItemChanged(item)"
            [(clrChecked)]="item.checked"
            [clrInline]="true"
            [clrDisabled]="item.disabled">
            <span>{{ item.value }}</span>
          </clr-checkbox>
        </li>
      </ul>
    </div>
  `
})

export class ShareFilterCheckboxComponent implements Filter<{ key: string, value: string }> {
  public ftype: string = FilterType.Checkbox;

  public static instanceof(obj: any) {
    return obj.hasOwnProperty('filterParamName') && obj.hasOwnProperty('items') && obj.hasOwnProperty('selectedItems');
  }

  @Input() public filterParamName: string;
  @Input() public items: Array<{ key: string, value: string, checked: boolean }>;
  public selectedItems: Array<{ key: string, value: string }> = [];
  public changes = new Subject<any>();
  @Input('defaultFilterValues')
  public set defaultFilterValues(newValues: string[]) {
    if (newValues) {
      this.selectedItems = [];
      if (this.items && Array.isArray(this.items)) {
        for (let item of this.items) {
          if (newValues && Array.isArray(newValues)) {
            for (let newValue of newValues) {
              if (item.key === newValue) {
                item.checked = true;
                this.selectedItems.push(item);
              }
            }
          }
        }
      }
    }
    this.changes.next(true);
  }

  public onItemChanged(item) {
    if (item.checked) {
      this.selectedItems.push(item);
    } else {
      let index = this.selectedItems.indexOf(item);
      if (index >= 0) {
        this.selectedItems.splice(index, 1);
      }
    }
    this.changes.next(true);
  }

  public accepts(item): boolean {
    for (let currentItem of this.items) {
      if (currentItem.checked && currentItem.key === item[this.filterParamName]) {
        return true;
      }
    }
    return false;
  }

  public isActive(): boolean {
    return this.selectedItems != null && this.selectedItems.length > 0;
  }
}
