import { Component, OnInit, ViewChild } from '@angular/core';

import { Item }                 from '../../../model/quant/item';
import { ItemService }          from '../../../service/quant/item.service';
import { ClassifyService }      from '../../../service/quant/classify.service';
import { AssetSource }          from '../../../model/quant/config';
import { ConfigAssetSourceComponent } from '../config/asset-source/asset-source.component';

@Component({
  selector: 'quant-item',
  templateUrl: './item.component.html',
  styleUrls: ['./item.component.css']
})
export class ItemComponent implements OnInit {
  @ViewChild(ConfigAssetSourceComponent)
  assetSourceComponent: ConfigAssetSourceComponent;

  items: Item[];
  allItems: Item[] = [];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  syncModelOpened: boolean = false;

  constructor(
    private itemService: ItemService,
    private classifyService: ClassifyService
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  Sync($event: any): void {
    if ($event) {
      this.classifyService.ListByAssetSource($event).subscribe(res => {
        if (res.length > 0) {
          this.syncModelOpened = false;
          this.itemService.WsSyncSource(res).subscribe(res => {
            if (res) {
              this.refresh();
            }
          })
        }
      });
    } else {
      this.syncModelOpened = false;
    }
  }

  openSyncModel(id?: number): void {
    this.syncModelOpened = true;
    this.assetSourceComponent.New()
  }

  delete(item: Item): void {
    this.itemService.Delete(item.Id).subscribe(res => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page) {
      if (this.allItems.length == 0) {
        this.refreshItem(state.page.from, state.page.to + 1);
      } else {
        this.items = this.allItems.slice(state.page.from, state.page.to + 1);
      }
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshItem(0, 10);
  }

  refreshItem(from: number, to: number): void {
    this.itemService.List(null).subscribe(res => {
      this.totalCount = res.length;
      this.allItems = res;
      this.items = res.slice(from, to);
    })
  }

  listClassifyName(item: Item): string {
    let ret: string = '';
    item.Classify.forEach((classify, key) => {
      ret += classify.Name + "\r\n";
    });
    return ret;
  }
}
