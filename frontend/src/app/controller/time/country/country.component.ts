import { Component, OnInit, ViewChild } from '@angular/core';

import { Country }             from '../../../model/time/country';
import { CountryService }      from '../../../service/time/country.service';
import { SaveCountryComponent }            from './save/save.component';

@Component({
  selector: 'time-country',
  templateUrl: './country.component.html',
  styleUrls: ['./country.component.css']
})
export class CountryComponent implements OnInit {
  @ViewChild(SaveCountryComponent)
  saveCountry: SaveCountryComponent;

  countries: Country[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private countryService: CountryService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    this.saveCountry.New(id);
  }

  delete(country: Country): void {
    this.countryService.Delete(country.Id).subscribe(res => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page) {
      this.refreshClassify(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshClassify(0, 10);
  }

  refreshClassify(from: number, to: number): void {
    this.countryService.List().subscribe(res => {
      this.totalCount = res.length;
      this.countries = res.slice(from, to);
    })
  }
}
