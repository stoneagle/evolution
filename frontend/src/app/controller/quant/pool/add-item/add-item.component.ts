import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Item }        from '../../../../model/quant/item';
import { Pool }        from '../../../../model/quant/pool';
import { Classify }    from '../../../../model/quant/classify';
import { PoolService } from '../../../../service/quant/pool.service';
import { ItemService } from '../../../../service/quant/item.service';
import { doFiltering } from '../../../../shared/utils';

@Component({
  selector: 'quant-pool-add-item',
  templateUrl: './add-item.component.html',
  styleUrls: ['./add-item.component.css']
})

export class PoolAddItemComponent implements OnInit {
  modelOpened: boolean = false;

  pool: Pool = new Pool();

  showItems: Item[]         = [];
  showAllItems: Item[]      = [];
  selectedShowItems: Item[] = [];

  submitItems: Item[]         = [];
  submitAllItems: Item[]      = [];
  selectedSubmitItems: Item[] = [];

  @Output() add = new EventEmitter<boolean>();

  showPageSize: number    = 5;
  showTotalCount: number  = 0;
  showCurrentPage: number = 1;

  submitPageSize: number    = 5;
  submitTotalCount: number  = 0;
  submitCurrentPage: number = 1;

  constructor(
    private poolService: PoolService,
    private itemService: ItemService,
  ) { }

  ngOnInit() {
    this.initItemArray();
  }

  initItemArray(): void {
    this.showItems           = [];
    this.showAllItems        = [];
    this.selectedShowItems   = [];
    this.submitItems         = [];
    this.submitAllItems      = [];
    this.selectedSubmitItems = [];
  }

  showLoad(state: any): void {
    if (state && state.page) {
      let filterShowItems = doFiltering<Item>(this.showAllItems, state);
      this.showItems = filterShowItems.slice(state.page.from, state.page.to + 1);
      this.showTotalCount = filterShowItems.length;
    }
  }

  submitLoad(state: any): void {
    if (state && state.page) {
      this.submitItems = this.submitAllItems.slice(state.page.from, state.page.to + 1);
    }
  }

  submitRefresh(): void {
    this.submitItems = this.submitAllItems.slice(0, this.submitPageSize);
    this.submitCurrentPage = 1;
    this.submitTotalCount = this.submitAllItems.length;
  }

  New(pool: Pool): void {
    this.initItemArray();
    this.pool = pool;

    this.submitAllItems = this.pool.Item;
    this.submitItems = this.pool.Item;

    let item = new Item();
    item.Classify = []; 
    let classify = new Classify();
    classify.AssetString = this.pool.Asset
    classify.Type = this.pool.Type
    item.Classify.push(classify)
    this.itemService.ListWithCondition(item).subscribe(res => {
      this.showTotalCount = res.length;
      this.showAllItems = res;
      this.showItems = res.slice(0, this.showPageSize);
      this.modelOpened = true;
    })
  }            

  Submit(): void {
    if (this.submitAllItems.length > 0) {
      this.pool.Item = this.submitAllItems;
      this.poolService.AddItems(this.pool).subscribe(res => {
        if (res) {
          this.add.emit(true);
          this.modelOpened = false;
        }
        this.pool.Item = [];
      })
    }
  }

  onAdd(): void {
    this.selectedShowItems.forEach((item, key) => {
      for (var i = 0; i < this.submitAllItems.length; i++) {
        if (item === this.submitAllItems[i]) {
          return;
        }
      }
      this.submitAllItems.push(item);
    });
    this.selectedShowItems = [];
    this.submitRefresh();
  }

  onDelete(): void {
    let tmpSubmitAllItems = [];
    this.submitAllItems.forEach((item, key) => {
      let tmpFlag = false;
      for (var i = 0; i < this.selectedSubmitItems.length; i++) {
        if (item === this.selectedSubmitItems[i]) {
          tmpFlag = true;
          continue;
        }
      }
      if (!tmpFlag) {
        tmpSubmitAllItems.push(item);
      }
    });
    this.selectedSubmitItems = [];
    this.submitAllItems = tmpSubmitAllItems;
    this.submitRefresh();
  }

  selectedShowItemsChange(): void {
  }

  selectedSubmitItemsChange(): void {
  }
}
