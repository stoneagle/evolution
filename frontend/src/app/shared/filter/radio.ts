import { Component, Input } from '@angular/core';
import { Filter }           from 'clarity-angular';
import { Subject }          from 'rxjs/Subject';
import { FilterType }       from '../const';

@Component({
  moduleId: 'share-filter-radio',
  selector: 'share-filter-radio',
  template: `
    <div>
      <form name="radioForm">
        <ul style="list-style: none">
          <li *ngFor="let item of items">
            <input type="radio" name="selectedItemName" [value]="item"
                   [(ngModel)]="selectedItem" (change)="onItemChanged(item)">
            <label for="selectedItemName">{{ item.value }}</label>
          </li>
        </ul>
      </form>
    </div>
  `
})
export class ShareFilterRadioComponent implements Filter<{ key: any, value: string }> {
  public ftype: string = FilterType.Radio;

  public static instanceof(obj: any) {
    return obj.hasOwnProperty('filterParamName') && obj.hasOwnProperty('items') && obj.hasOwnProperty('selectedItem');
  }

  @Input() public filterParamName: string;
  @Input() public items: Array<{ key: any, value: string }>;
  @Input('defaultFilterValue')
  public set defaultFilterValue(newValue: string) {
    for (let item of this.items) {
      if (item.key === newValue) {
        this.selectedItem = item;
        this.onItemChanged(item);
      }
    }
  }

  public selectedItem: { key: any, value: string };
  public changes = new Subject<any>();
  public onItemChanged(item) {
    this.changes.next(true);
  }

  public accepts(item): boolean {
    return item[this.filterParamName] === this.selectedItem.key ;
  }

  public isActive(): boolean {
    return this.selectedItem != null && this.selectedItem !== this.items[0];
  }
}

