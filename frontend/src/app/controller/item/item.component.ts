import { Component, OnInit, ViewChild } from '@angular/core';
import { Item } from '../../model/business/item';
import { ItemService } from '../../service/business/item.service';
import { SyncItemComponent } from './sync/sync.component';

@Component({
  selector: 'app-item',
  templateUrl: './item.component.html',
  styleUrls: ['./item.component.css']
})
export class ItemComponent implements OnInit {
  @ViewChild(SyncItemComponent)
  syncItem: SyncItemComponent;

  items: Item[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private itemService: ItemService
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  syncd(syncd: boolean): void {
    if (syncd) {
      this.refresh();
    }
  }

  openSyncModel(id?: number): void {
    this.syncItem.New(id);
  }

  delete(item: Item): void {
    this.itemService.Delete(item.Id)
    .subscribe(res => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page) {
      this.refreshItem(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshItem(0, 10);
  }

  refreshItem(from: number, to: number): void {
    this.itemService.List()
    .subscribe(res => {
      this.totalCount = res.length;
      this.items = res.slice(from, to);
    })
  }
}
